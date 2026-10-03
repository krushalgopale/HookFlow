package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
)

func SaveAPIKey(
	db *pgxpool.Pool,
	ctx context.Context,
	apikeyID string,
	environmentID string,
	key string,
	name string,
	expiresAt *time.Time,
) error {

	hash := sha256.Sum256([]byte (key))
	keyHash := hex.EncodeToString(hash[:])

	_, err := db.Exec(
		ctx,
		`INSERT INTO api_keys(id, environment_id, name, key, expires_at) VALUES ($1, $2, $3, $4, $5)`,
		apikeyID,
		environmentID,
		name,
		keyHash,
		expiresAt,
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

func GetAPIKeyByID(
	db *pgxpool.Pool,
	ctx context.Context,
	apiKeyID string,
	environmentID string,
) (models.APIKey, error) {
	var apiKey models.APIKey

	err := db.QueryRow(
		ctx,
		`SELECT id, environment_id, name, created_at, updated_at, expires_at FROM api_keys WHERE id = $1 AND environment_id = $2`,
		apiKeyID,
		environmentID,
	).Scan(
		&apiKey.ID,
		&apiKey.EnvironmentID,
		&apiKey.Name,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
		&apiKey.ExpiresAt,
	)
	if err != nil {
		return models.APIKey{}, err
	}

	return apiKey, nil
}

func ListAPIKeysByEnvironment(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
) ([]models.APIKey, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, environment_id, name, created_at, updated_at, expires_at FROM api_keys WHERE environment_id = $1`,
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
			&apikey.Name,
			&apikey.CreatedAt,
			&apikey.UpdatedAt,
			&apikey.ExpiresAt,
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
