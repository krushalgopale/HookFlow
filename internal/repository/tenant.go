package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
	tenantName string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO tenants (id, name)
		VALUES ($1, $2)`,
		tenantID,
		tenantName,
	)
	return err
}
