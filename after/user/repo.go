package user

import (
	"context"
	"database/sql"
	"fmt"
)

const saveQuery string = "INSERT INTO ..."

type Repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, u User) error {
	_, err := r.db.ExecContext(ctx, saveQuery, u.ID, u.Email, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to persist user: %w", err)
	}

	return nil
}
