package db

import (
	"database/sql"
)

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
