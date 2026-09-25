package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
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

func ListTenants(
	db *pgxpool.Pool,
	ctx context.Context,
) ([]models.Tenant, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, name, created_at, updated_at FROM tenants ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []models.Tenant

	for rows.Next() {
		var tenant models.Tenant

		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tenants = append(tenants, tenant)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}
