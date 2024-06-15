package dbrepo

type DB interface {
	Connect(string) error
	Close() error
	Migrate() error
}
