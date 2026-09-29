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
	userID string,
	tenantName string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO tenants (id, user_id, name)
		VALUES ($1, $2, $3)`,
		tenantID,
		userID,
		tenantName,
	)
	return err
}

func GetTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
	userID string,
) (*models.Tenant, error) {
	var tenant models.Tenant

	err := db.QueryRow(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM tenants WHERE id = $1 AND user_id = $2`,
		tenantID,
		userID,
	).Scan(&tenant.ID, &tenant.UserID, &tenant.Name, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		return &models.Tenant{}, err
	}

	return &tenant, nil
}

func ListTenants(
	db *pgxpool.Pool,
	ctx context.Context,
	userID string,
) ([]models.Tenant, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM tenants WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
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
			&tenant.UserID,
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
	userID string,
	tenantName string,
) error {
	result, err := db.Exec(
		ctx,
		`UPDATE tenants SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND user_id = $3`,
		tenantName,
		tenantID,
		userID,
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
	userID string,
) error {
	result, err := db.Exec(
		ctx,
		`DELETE FROM tenants WHERE id = $1 AND user_id = $2`,
		tenantID,
		userID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
