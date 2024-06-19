package dbrepo

// DB defines an interface for the database.
type DB interface {
	Connect(dsn string) error
	Close() error
	Migrate() error
}
