package env

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
	"github.com/luanapp/gin-example/pkg/logger"
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	sugar = logger.New()

	envFiles := []string{".env"}
	currEnv := os.Getenv("USING_GIN_ENV")
	sugar.Infof("curent environment: %s", currEnv)

	if currEnv == "" {
		envFiles = append(envFiles, ".env.local")
	} else {
		envFiles = append(envFiles, fmt.Sprintf("%s/%s", ".env.", currEnv))
	}

	err := godotenv.Load(envFiles...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
