package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
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
	Instance   = new(Env)
	production = map[string]string{"production": "production", "prod": "prod"}
)

func init() {
	loadEnvFile()

	Instance.Environment = getEnvOrDefault(environment, "development")
	Instance.Server.Port = getEnvOrDefault(port, "development")
	Instance.Database.URL = getEnvOrDefault(databaseURL, "")
	Instance.Database.Schema = getEnvOrDefault(databaseSchema, "")
}

func loadEnvFile() {
	envFiles := []string{".env"}

	_ = godotenv.Load(envFiles...)
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
