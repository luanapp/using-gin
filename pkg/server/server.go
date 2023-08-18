package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/luanapp/gin-example/config/env"
	"github.com/luanapp/gin-example/pkg/crud"
	"github.com/luanapp/gin-example/pkg/logger"
	"github.com/luanapp/gin-example/pkg/model"
	_ "github.com/luanapp/gin-example/pkg/server/docs"
	"github.com/penglongli/gin-metrics/ginmetrics"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"go.uber.org/zap"
)

type Server struct {
}

var (
	sugar *zap.SugaredLogger
)

func init() {
	sugar = logger.New()
}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Start() {
	r := setupEngine()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", env.Instance.Server.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			sugar.Fatalf("listen: %s\n", err)
		}
		sugar.Infof("server is up in the port %d", env.Instance.Server.Port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		sugar.Fatalw("Server forced to shutdown...", "error", err.Error())
	}

}

func setupEngine() *gin.Engine {
	r := gin.New()
	setupMiddlewares(r)
	setupMetrics(r)
	addHelperRoutes(r)
	addCrudRoutes[model.Species](r, "/species")
	addCrudRoutes[model.Sample](r, "/sample")
	return r
}

func setupMiddlewares(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{http.MethodHead, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:    []string{"*"},
	}))
	l := logger.NewLogger()
	r.Use(ginZap.Ginzap(l, time.RFC3339, true))
	r.Use(ginZap.RecoveryWithZap(l, true))
}

func setupMetrics(r *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetMetricPrefix("nhm_")
	m.Use(r)
}

func addHelperRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthHandler := NewHandler()
	healthRoute := r.Group("/status")
	healthRoute.GET("", gin.WrapF(healthHandler.StatusHandler()))
	healthRoute.GET("/health", healthHandler.Health)
}

func addCrudRoutes[T any](r *gin.Engine, relativePath string) {
	handler := crud.DefaultHandler[T]()
	route := r.Group(relativePath)
	route.GET("", handler.GetAll)
	route.GET("/:id", handler.GetById)
	route.POST("", handler.Save)
	route.PUT("/:id", handler.Update)
	route.DELETE("/:id", handler.Delete)
}
