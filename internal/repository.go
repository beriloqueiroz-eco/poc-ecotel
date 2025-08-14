package internal

import (
	"context"
)

type Repository struct {
	DB *DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) SimpleQuery(ctx context.Context) (string, error) {
	row := r.DB.QueryRow(ctx, "SELECT cod from country")
	var result string
	if err := row.Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}
