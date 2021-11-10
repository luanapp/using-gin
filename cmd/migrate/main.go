package main

import (
	"os"

	"go.uber.org/zap"
	"luana.com/gin-example/config/database"
	"luana.com/gin-example/pkg/logger"
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	sugar = logger.New()
}

func main() {
	database.InitializeDB()

	if len(os.Args) == 1 {
		filename := os.Args[1]
		if database.MigrationExists(filename) {
			sugar.Warn("migration was already applied")
			return
		}
		err := database.UpFromFilename(filename)
		if err != nil {
			sugar.Fatalw("there was a problem when applying migration", "error", err.Error())
		}
	} else {
		err := database.Up()
		if err != nil {
			sugar.Fatalw("there was a problem when applying migration", "error", err.Error())
		}
	}
}
