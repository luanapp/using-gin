package database

import (
	"context"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/luanapp/gin-example/pkg/logger"
	"go.uber.org/zap"
)

const createSchemaSQL = "CREATE SCHEMA IF NOT EXISTS natural_history_museum;"

var (
	connection *pgxpool.Pool
	sugar      *zap.SugaredLogger
)

func init() {
	sugar = logger.New()
}

// InitializeDB this will start a connection pool with the database with the url DATABASE_URL
// It will create the application schema if not exists. The same with the migration table
func InitializeDB() {
	var err error
	url := os.Getenv("DATABASE_URL")
	sugar.Infof("connecting to database %s", strings.Split(url, "@")[1])
	connection, err = pgxpool.Connect(context.Background(), url)
	if err != nil {
		sugar.Fatal(err)
	}

	tx, err := connection.Begin(context.Background())
	if err != nil {
		sugar.Fatal()
	}

	err = createSchema(tx)
	if err != nil {
		rollback(tx)
		sugar.Fatal(err)
	}

	err = createMigrationTable(tx)
	if err != nil {
		rollback(tx)
		sugar.Fatal(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		sugar.Fatal(err)
	}
}

// GetConnection retrieve the application connection pool with the database
func GetConnection() *pgxpool.Pool {
	return connection
}

func rollback(tx pgx.Tx) {
	err := tx.Rollback(context.Background())
	if err != nil {
		sugar.Fatal(err)
	}
}

func createSchema(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), createSchemaSQL)
	if err != nil {
		return err
	}
	sugar.Info("natural_history_museum schema was created or already exists")
	return nil
}
