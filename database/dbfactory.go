package database

type DB interface {
	NewDB() DataStore
}
