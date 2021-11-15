package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/luanapp/gin-example/pkg/domain/species"
	_ "github.com/luanapp/gin-example/pkg/env"
	"github.com/luanapp/gin-example/pkg/logger"
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

	spHandler := species.NewHandler(species.DefaultRepository())
	spRoute := r.Group("/species")
	spRoute.GET("", spHandler.GetAll)
	spRoute.GET("/:id", spHandler.GetById)
	spRoute.POST("", spHandler.Save)
	spRoute.PUT("/:id", spHandler.Update)
	spRoute.DELETE("/:id", spHandler.Delete)

	srv := http.Server{
		Addr:    ":8080",
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
