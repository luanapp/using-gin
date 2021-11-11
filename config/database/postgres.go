package database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/luanapp/gin-example/pkg/logger"
)

const createSchemaSQL = "CREATE SCHEMA IF NOT EXISTS natural_history_museum;"

var (
	connection *pgx.Conn
	sugar      *zap.SugaredLogger
)

func init() {
	sugar = logger.New()
}

func InitializeDB() {
	var err error
	url := "postgres://postgres:postgres@localhost:5432/postgres"
	connection, err = pgx.Connect(context.Background(), url)
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

func GetConnection() *pgx.Conn {
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
	sugar.Info("created natural_history_museum schema successfully")
	return nil
}
