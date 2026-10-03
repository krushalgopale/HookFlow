package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveAPIKey(
	db *pgxpool.Pool,
	ctx context.Context,
	apikeyID string,
	environmentID string,
	key string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO api_keys(id, environment_id, key) VALUES ($1, $2, $3)`,
		apikeyID,
		environmentID,
		key,
	)

	return err
}

func EnvironmentBelongsToUser(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
	userID string,
) (bool, error) {
	var exists bool

	err := db.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM environments e INNER JOIN tenants t ON e.tenant_id = t.id WHERE e.id = $1 AND t.user_id = $2)`,
		environmentID,
		userID,
	).Scan(&exists)

	return exists, err
}
