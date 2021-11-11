package database

import (
	"context"

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

func InitializeDB() {
	var err error
	url := "postgres://postgres:postgres@localhost:5432/postgres"
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
