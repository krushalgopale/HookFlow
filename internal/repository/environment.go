package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
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

func GetEnvironmentByTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
	tenantID string,
) (models.Environment, error) {
	var environment models.Environment

	err := db.QueryRow(
		ctx,
		`SELECT id, tenant_id, name, created_at, updated_at FROM environments WHERE id = $1 AND tenant_id = $2`,
		environmentID,
		tenantID,
	).Scan(
		&environment.ID,
		&environment.TenantID,
		&environment.Name,
		&environment.CreatedAt,
		&environment.UpdatedAt,
	)
	if err != nil {
		return models.Environment{}, err
	}

	return environment, nil
}

func ListEnvironmentsByTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
) ([]models.Environment, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, tenant_id, name, created_at, updated_at FROM environments WHERE tenant_id = $1`,
		tenantID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var environments []models.Environment

	for rows.Next() {
		var environment models.Environment

		err := rows.Scan(
			&environment.ID,
			&environment.TenantID,
			&environment.Name,
			&environment.CreatedAt,
			&environment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		environments = append(environments, environment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return environments, nil
}

func UpdateEnvironmentByTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	envronmentID string,
	tenantID string,
	name string,
) error {
	result, err := db.Exec(
		ctx,
		`UPDATE environments SET NAME = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND 	tenant_id = $3`,
		name,
		envronmentID,
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

func DeleteEnvironmentByTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
	tenantID string,
) error {
	result, err := db.Exec(
		ctx,
		`DELETE FROM environments WHERE id = $1 AND tenant_id = $2`,
		environmentID,
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
