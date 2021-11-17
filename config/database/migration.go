package database

import (
	"context"
	"fmt"
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
	selectAllMigrationsSQL = "SELECT m.file FROM natural_history_museum.migration m;"
	searchMigrationSQL     = "SELECT m.id FROM natural_history_museum.migration m WHERE m.file = $1;"
	insertMigrationSQL     = "INSERT INTO natural_history_museum.migration (file) VALUES ($1);"
	deleteMigrationSQL     = "DELETE FROM natural_history_museum.migration m WHERE m.file = $1 ;"
)

func createMigrationTable(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), createMigrationTableSQL)
	if err != nil {
		return err
	}
	sugar.Info("migration table was created or already exists")
	return nil
}

// UpFromFilename runs the up migration in the given file, if not yet applied
func UpFromFilename(filename string) error {
	file, err := os.Lstat(fmt.Sprintf("%s/%s", migrationsBaseDir, filename))
	if err != nil {
		sugar.Errorw("failed to apply migration", "file", file, "error", err.Error())
		return err
	}

	m := readMigrationFromFile(migrationsBaseDir, file.Name())
	return upMigration(m)
}

// Up run all migrations not yet applied in the database
func Up() error {
	migrations := readMigrationFromFiles()

	for _, m := range migrations {
		err := upMigration(m)
		if err != nil {
			sugar.Errorw("failed to apply migration", "file", m.filename, "error", err.Error())
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

	if !migrationExists(m.filename) {
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

// DownFromFilename runs the down script in the database in the given filename
func DownFromFilename(filename string) error {
	file, err := os.Lstat(fmt.Sprintf("%s/%s", migrationsBaseDir, filename))
	if err != nil {
		sugar.Errorw("failed to undo migration", "file", file, "error", err.Error())
		return err
	}

	m := readMigrationFromFile(migrationsBaseDir, file.Name())
	return downMigration(m)
}

// Down run all down scripts applied in the database
func Down() error {
	dbFiles := migrationFilesFromDB()

	for _, file := range dbFiles {
		m := readMigrationFromFile(migrationsBaseDir, file)
		err := downMigration(m)
		if err != nil {
			sugar.Errorw("failed to undo migration", "file", file, "error", err.Error())
			return err
		}
	}
	return nil
}

func downMigration(m *migration) error {
	tx, err := prepareMigrationTransaction()
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), m.Down)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}
	_, err = tx.Exec(context.Background(), deleteMigrationSQL, m.filename)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}

	_ = tx.Commit(context.Background())
	return nil
}

/*func generateKey() string {
	return time.Now().Format("20060102150405")
}*/

func migrationExists(filename string) bool {
	result, err := connection.Query(context.Background(), searchMigrationSQL, filename)
	if err != nil {
		return false
	}
	return result.Next()
}

func prepareMigrationTransaction() (pgx.Tx, error) {
	txOpt := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	return connection.BeginTx(context.Background(), txOpt)
}

func migrationFilesFromDB() []string {
	rows, err := connection.Query(context.Background(), selectAllMigrationsSQL)
	if err != nil {
		return nil
	}

	var files []string
	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err != nil {
			return nil
		}
		files = append(files, file)
	}
	return files
}

func readMigrationFromFiles() []*migration {
	var mFiles []*migration

	files, err := ioutil.ReadDir(migrationsBaseDir)
	if err != nil {
		sugar.Fatal(err)
	}

	for _, file := range files {
		mFile := readMigrationFromFile(migrationsBaseDir, file.Name())
		mFiles = append(mFiles, mFile)
	}
	return mFiles
}

func readMigrationFromFile(baseDir string, file string) *migration {
	mFile := &migration{}
	fullPath := fmt.Sprintf("%s/%s", baseDir, file)
	bytes, _ := ioutil.ReadFile(fullPath)
	err := yaml.Unmarshal(bytes, mFile)
	if err != nil {
		sugar.Fatal(err)
	}
	mFile.filename = file
	return mFile
}
