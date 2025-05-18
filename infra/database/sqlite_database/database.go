package sqlite_database

import (
	"database/sql"
	"os"
	"path/filepath"

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
	if err != nil {
		panic(err)
	}
	dbFilePath := filepath.Join(configPath, ".thop", DB_FILE_NAME)

	db, err := sql.Open(DB_DRIVER, dbFilePath)
	if err != nil {
		panic(err)
	}
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
	if err != nil {
		panic(err)
	}

	// validate migrations
	var migrationsCount int
	err = db.QueryRow(sql_count_migrations).Scan(&migrationsCount)
	if err != nil {
		panic(err)
	}

	// run migrations if needed
	for i := migrationsCount; i < len(migrations); i++ {
		migration := migrations[i]

		// exec migration
		_, err := db.Exec(migration.sql)
		if err != nil {
			panic(err)
		}

		// save execution in the database
		_, err = db.Exec(sql_insert_migration, i, migration.description)
		if err != nil {
			panic(err)
		}
	}
}
