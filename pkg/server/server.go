package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/luanapp/gin-example/pkg/crud"
	_ "github.com/luanapp/gin-example/pkg/env"
	"github.com/luanapp/gin-example/pkg/logger"
	"github.com/luanapp/gin-example/pkg/model"
	_ "github.com/luanapp/gin-example/pkg/server/docs"

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
	return &Server{}
}

func (s *Server) Start() {
	r := setupEngine()
	setupMiddlewares(r)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("listen: %s\n", err)
		}
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

func setupMiddlewares(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{http.MethodHead, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:    []string{"*"},
	}))
}

func setupEngine() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthHandler := NewHandler()
	healthRoute := r.Group("/status")
	healthRoute.GET("", gin.WrapF(healthHandler.StatusHandler()))
	healthRoute.GET("/health", healthHandler.Health)

	addCrudRoutes[model.Species](r, "/species")
	addCrudRoutes[model.Sample](r, "/sample")
	return r
}

func addCrudRoutes[T model.Model](r *gin.Engine, relativePath string) {
	handler := crud.DefaultHandler[T]()
	route := r.Group(relativePath)
	route.GET("", handler.GetAll)
	route.GET("/:id", handler.GetById)
	route.POST("", handler.Save)
	route.PUT("/:id", handler.Update)
	route.DELETE("/:id", handler.Delete)
}
