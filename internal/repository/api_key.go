package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
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

func ListAPIKeysByEnvironment(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
) ([]models.APIKey, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, environment_id, key, created_at FROM api_keys WHERE environment_id = $1`,
		environmentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []models.APIKey

	for rows.Next() {
		var apikey models.APIKey

		err := rows.Scan(
			&apikey.ID,
			&apikey.EnvironmentID,
			&apikey.Key,
			&apikey.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		apiKeys = append(apiKeys, apikey)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}
	return apiKeys, nil
}
