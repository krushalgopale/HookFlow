package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveUser(
	db *pgxpool.Pool,
	ctx context.Context,
	userID string,
	email string,
	passwordHash string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`,
		userID,
		email,
		passwordHash,
	)

	return err
}
