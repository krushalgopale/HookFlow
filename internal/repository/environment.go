package repository

import (
	"context"

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

func ListEnvironmentsByTenant(
	db *pgxpool.Pool,
	ctx context.Context,
	tenantID string,
) ([]models.Environment, error){
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
