package database

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	"go.uber.org/zap"

	"github.com/luanapp/gin-example/pkg/env"
	"github.com/luanapp/gin-example/pkg/logger"
)

const createSchemaSQL = "CREATE SCHEMA IF NOT EXISTS natural_history_museum;"

var (
	connection *pgxpool.Pool
	sugar      *zap.SugaredLogger
)

func init() {
	sugar = logger.New()
}

// GetConnection retrieve the application connection pool with the database
func GetConnection() *pgxpool.Pool {
	return connection
}

// InitializeDB this will start a connection pool with the database with the url DATABASE_URL
// It will create the application schema if not exists. The same with the migration table
func InitializeDB() {
	var err error

	url := env.Instance.Database.URL
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		sugar.Fatal(err)
	}

	initializeConfig(config)

	sugar.Infof("connecting to database %s", strings.Split(url, "@")[1])
	connection, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		sugar.Fatal(err)
	}

	initialDatabaseData()
}

func initializeConfig(config *pgxpool.Config) {
	config.ConnConfig.LogLevel = pgx.LogLevelInfo
	config.ConnConfig.Logger = zapadapter.NewLogger(logger.NewLogger())
	config.ConnConfig.RuntimeParams = map[string]string{
		"search_path": env.Instance.Database.Schema,
	}
}

func initialDatabaseData() {
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
