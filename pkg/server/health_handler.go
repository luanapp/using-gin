package server

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

// Health server health check endpoint
// @Summary Returns OK if the server is up
// @Description This returns a JSON {"status": "OK"} with a status 200 if the server is up
// @Tags health
// @Success 200 {object} object
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *Handler) StatusHandler() http.HandlerFunc {
	statusHealth, _ := healthgo.New(healthgo.WithChecks(
		healthgo.Config{
			Name:    "server-status",
			Timeout: time.Second * 2,
			Check: httpCheck.New(httpCheck.Config{
				URL: fmt.Sprintf("http://localhost:%s/status/health", os.Getenv("PORT")),
			})},
		healthgo.Config{
			Name:    "postgres-status",
			Timeout: time.Second * 3,
			Check: postgresCheck.New(postgresCheck.Config{
				DSN: fmt.Sprintf("%s?%s", os.Getenv("DATABASE_URL"), "sslmode=disable"),
			}),
		},
	))
	return statusHealth.HandlerFunc
}
