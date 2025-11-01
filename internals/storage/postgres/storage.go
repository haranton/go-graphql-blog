package postgres

import (
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (st *PostgresStorage) Close() error {
	return st.db.Close()
}
