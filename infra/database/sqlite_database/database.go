package sqlite_database

import (
	"database/sql"
	"os"
	"path/filepath"
	"phopper/domain/errors"

	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE_NAME = "thop.db"
const DB_DRIVER = "sqlite3"
const sql_create_migrations_table = `CREATE TABLE IF NOT EXISTS migrations (version INTEGER PRIMARY KEY, description TEXT);`
const sql_insert_migration = `INSERT INTO migrations (version, description) VALUES (?, ?);`
const sql_count_migrations = `SELECT COUNT(*) FROM migrations;`

type SqliteDatabase struct {}

// util function to not retype driver and file path every time
func openDB() *sql.DB {
	configPath, err := os.UserConfigDir()
	errors.EnsureNotNil(err, "Could not get user config directory")

	dbFilePath := filepath.Join(configPath, ".thop", DB_FILE_NAME)

	db, err := sql.Open(DB_DRIVER, dbFilePath)
	errors.EnsureNotNil(err, "Could not open database")

	return db
}

func (_ SqliteDatabase) GetProjectRepository() SqliteProjectRepository {
	return SqliteProjectRepository{}
}

func (_ SqliteDatabase) RunMigrations() {
	migrations := getMigrations()
	
	db := openDB()
	defer db.Close()

	// create migrations table if not exists
	_, err := db.Exec(sql_create_migrations_table)
	errors.EnsureNotNil(err, "Could not create migrations table")

	// validate migrations
	var migrationsCount int
	err = db.QueryRow(sql_count_migrations).Scan(&migrationsCount)
	errors.EnsureNotNil(err, "Could not count migrations")

	// run migrations if needed
	for i := migrationsCount; i < len(migrations); i++ {
		migration := migrations[i]

		// exec migration
		_, err := db.Exec(migration.sql)
		errors.EnsureNotNil(err, "Could not run migration")

		// save migration run
		_, err = db.Exec(sql_insert_migration, i, migration.description)
		errors.EnsureNotNil(err, "Could not save migration run")
	}
}
