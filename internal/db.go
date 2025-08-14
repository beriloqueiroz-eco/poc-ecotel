package internal

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
	Builder *squirrel.StatementBuilderType
}
type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
}

func NewDB(ctx context.Context, cfg DBConfig, pool *pgxpool.Pool) (*DB, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &DB{
		Pool:    pool,
		Builder: &builder,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
