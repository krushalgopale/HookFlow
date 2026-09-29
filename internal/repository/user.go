package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
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

func GetUserByEmail(
	db *pgxpool.Pool,
	ctx context.Context,
	email string,
) (models.User, error) {
	var user models.User

	err := db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByID(
	db *pgxpool.Pool,
	ctx context.Context,
	userID string,
) (models.User, error) {
	var user models.User

	err := db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
