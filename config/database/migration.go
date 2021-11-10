package database

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/jackc/pgx/v4"
	"gopkg.in/yaml.v3"
)

type (
	migration struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Up          string `yaml:"up"`
		Down        string `yaml:"down"`
		filename    string
	}
)

const (
	migrationsBaseDir       = "./migrations"
	createMigrationTableSQL = `
CREATE TABLE IF NOT EXISTS natural_history_museum.migration (
    id SERIAL CONSTRAINT migration_id_key PRIMARY KEY,
    file varchar(300) CONSTRAINT migration_file_unique UNIQUE,
    date timestamp DEFAULT now()
);
`
	searchMigrationSQL = "SELECT m.id FROM natural_history_museum.migration m WHERE m.file = $1"
	insertMigrationSQL = "INSERT INTO natural_history_museum.migration (file) VALUES ($1);"
)

func createMigrationTable(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), createMigrationTableSQL)
	if err != nil {
		return err
	}
	sugar.Info("created migrationEntry table successfully")
	return nil
}

func UpFromFilename(filename string) error {
	file, err := os.Lstat(fmt.Sprintf("%s/%s", migrationsBaseDir, filename))
	if err != nil {
		return err
	}

	m := readMigrationFromFile(migrationsBaseDir, file)
	return upMigration(m)
}

func Up() error {
	migrations := readMigrationFromFiles()

	for _, m := range migrations {
		err := upMigration(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func upMigration(m *migration) error {
	tx, err := prepareMigrationTransaction()
	if err != nil {
		return err
	}

	if !MigrationExists(m.filename) {
		_, err = tx.Exec(context.Background(), m.Up)
		if err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}
		_, err = tx.Exec(context.Background(), insertMigrationSQL, m.filename)
		if err != nil {
			_ = tx.Rollback(context.Background())
			return err
		}
	} else {
		sugar.Warnf("migration %s already exists", m.filename)
	}
	_ = tx.Commit(context.Background())
	return nil
}

func prepareMigrationTransaction() (pgx.Tx, error) {
	txOpt := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	return connection.BeginTx(context.Background(), txOpt)
}

/*func (m migration) down() error {
	return nil
}*/

/*func generateKey() string {
	return time.Now().Format("20060102150405")
}*/

func MigrationExists(filename string) bool {
	result, err := connection.Query(context.Background(), searchMigrationSQL, filename)
	if err != nil {
		return false
	}
	return result.Next()
}

func readMigrationFromFiles() []*migration {
	var mFiles []*migration

	files, err := ioutil.ReadDir(migrationsBaseDir)
	if err != nil {
		sugar.Fatal(err)
	}

	for _, file := range files {
		mFile := readMigrationFromFile(migrationsBaseDir, file)
		mFiles = append(mFiles, mFile)
	}
	return mFiles
}

func readMigrationFromFile(baseDir string, file fs.FileInfo) *migration {
	mFile := &migration{}
	fullPath := fmt.Sprintf("%s/%s", baseDir, file.Name())
	bytes, _ := ioutil.ReadFile(fullPath)
	err := yaml.Unmarshal(bytes, mFile)
	if err != nil {
		sugar.Fatal(err)
	}
	mFile.filename = file.Name()
	return mFile
}
