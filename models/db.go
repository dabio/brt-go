package models

import "database/sql"

// Datastore implements methods used in main.
type Datastore interface {
	AllEvents(year int) ([]*Event, error)
}

// DB struct contains attribute for database connection.
type DB struct {
	*sql.DB
}

// NewDB creates and returns a newly created database struct.
func NewDB(driverName string, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
