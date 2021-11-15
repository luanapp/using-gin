package health

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	healthgo "github.com/hellofresh/health-go/v4"
	httpCheck "github.com/hellofresh/health-go/v4/checks/http"
	postgresCheck "github.com/hellofresh/health-go/v4/checks/postgres"
)

type (
	Handlerer interface {
		Health(c *gin.Context)
		StatusHandler() http.HandlerFunc
	}

	Handler struct {
	}
)

func NewHandler() Handlerer {
	return &Handler{}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *Handler) StatusHandler() http.HandlerFunc {
	statusHealth, _ := healthgo.New(healthgo.WithChecks(
		healthgo.Config{
			Name:    "server-status",
			Timeout: time.Second * 5,
			Check: httpCheck.New(httpCheck.Config{
				URL: fmt.Sprintf("http://localhost:%s/status/health", os.Getenv("PORT")),
			})},
		healthgo.Config{
			Name:    "postgres-status",
			Timeout: time.Second * 5,
			Check: postgresCheck.New(postgresCheck.Config{
				DSN: fmt.Sprintf("%s?%s", os.Getenv("DATABASE_URL"), "sslmode=disable"),
			}),
		},
	))
	return statusHealth.HandlerFunc
}
