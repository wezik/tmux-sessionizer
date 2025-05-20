package storage

type Storage interface {
	RunMigrations()
}
