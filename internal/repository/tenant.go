package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func GetTenant(
	db *pgxpool.Pool, 
	ctx context.Context,
	tenantID string,
) (*models.Tenant, error) {
	var tenant models.Tenant

	err := db.QueryRow(
		ctx,
		`SELECT id, name, created_at, updated_at FROM tenants WHERE id = $1`,
		tenantID,
		).Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, &tenant.UpdatedAt)

	if err != nil {
		return &models.Tenant{}, err
	}

	return &tenant, nil
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

func UpdateTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
	tenantName string,
) error {
	result, err := db.Exec(
		ctx,
		`UPDATE tenants SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`,
		tenantName,
		tenantID,
		)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func DeleteTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
) error {

	result, err := db.Exec(
		ctx,
		`DELETE FROM tenants WHERE id = $1`,
		tenantID,
		)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	} 

	return nil
}
