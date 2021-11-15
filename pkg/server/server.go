package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luanapp/gin-example/pkg/domain/health"
	"github.com/luanapp/gin-example/pkg/domain/species"
	_ "github.com/luanapp/gin-example/pkg/env"
	"github.com/luanapp/gin-example/pkg/logger"
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
	r := gin.Default()

	healthHandler := health.NewHandler()
	healthRoute := r.Group("/status")
	healthRoute.GET("", gin.WrapF(healthHandler.StatusHandler()))
	healthRoute.GET("/health", healthHandler.Health)

	spHandler := species.NewHandler(species.DefaultRepository())
	spRoute := r.Group("/species")
	spRoute.GET("", spHandler.GetAll)
	spRoute.GET("/:id", spHandler.GetById)
	spRoute.POST("", spHandler.Save)
	spRoute.PUT("/:id", spHandler.Update)
	spRoute.DELETE("/:id", spHandler.Delete)

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
