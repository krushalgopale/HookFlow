package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveEnvironment(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
	tenantID string,
	name string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO environments (id, tenant_id, name) VALUES ($1, $2, $3)`,
		environmentID,
		tenantID,
		name,
	)
	return err
}

func TenantBelongsToUser(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
	userID string,
) (bool, error) {
	var exists bool

	err := db.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM tenants WHERE id = $1 AND user_id = $2)`,
		tenantID,
		userID,
	).Scan(&exists)

	return exists, err
}
