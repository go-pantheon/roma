package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXPool interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Close()
}

type DB struct {
	DB PGXPool
}

func NewPostgreSQLDB(pdb *pgxpool.Pool) *DB {
	return &DB{DB: pdb}
}
