package database

type Database interface {
	RunMigrations()
}
