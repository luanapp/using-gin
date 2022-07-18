package env

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/joho/godotenv"

	"github.com/luanapp/gin-example/pkg/logger"
)

type (
	Env struct {
		Environment string
		Server      struct {
			Port string
		}
		Database struct {
			URL    string
			Schema string
		}
	}
)

const (
	port           = "PORT"
	databaseURL    = "DATABASE_URL"
	databaseSchema = "DATABASE_SCHEMA"
	environment    = "USING_GIN_ENV"
)

var (
	sugar      *zap.SugaredLogger
	Instance   = &Env{}
	production = map[string]string{"production": "production", "prod": "prod"}
)

func init() {
	sugar = logger.New()

	loadEnvFile()

	Instance.Environment = getEnvOrDefault(environment, "development")
	Instance.Server.Port = getEnvOrDefault(port, "development")
	Instance.Database.URL = getEnvOrDefault(databaseURL, "")
	Instance.Database.Schema = getEnvOrDefault(databaseSchema, "")
}

func loadEnvFile() {
	envFiles := []string{".env"}
	currEnv := os.Getenv("USING_GIN_ENV")
	sugar.Infof("curent environment: %s", currEnv)

	if currEnv == "" {
		envFiles = append(envFiles, ".env.local")
	} else if currEnv != "production" {
		envFiles = append(envFiles, fmt.Sprintf("%s/%s", ".env.", currEnv))
	}

	err := godotenv.Load(envFiles...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getEnvOrDefault(key string, defaultValue string) string {
	if os.Getenv(key) == "" {
		return defaultValue
	}
	return os.Getenv(key)
}

func IsProduction() bool {
	_, ok := production[strings.ToLower(Instance.Environment)]
	return ok
}
