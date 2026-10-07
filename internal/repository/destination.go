package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
)

func SaveDestination(
	db *pgxpool.Pool,
	ctx context.Context,
	destinationID string,
	environmentID string,
	name string,
	url string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO destinations (id, environment_id, name, url) VALUES ($1, $2, $3, $4)`,
		destinationID,
		environmentID,
		name,
		url,
	)

	return err
}

func GetDestinationByID(
	db *pgxpool.Pool,
	ctx context.Context,
	destinationID string,
	environmentID string,
) (*models.Destination, error) {
	var destination models.Destination
	err := db.QueryRow(
		ctx,
		`SELECT id, environment_id, name, url, created_at, updated_at FROM destinations WHERE id = $1 AND environment_id = $2`,
		destinationID,
		environmentID,
	).Scan(
		&destination.ID,
		&destination.EnvironmentID,
		&destination.Name,
		&destination.URL,
		&destination.CreatedAt,
		&destination.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

func ListDestinations(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
) ([]models.Destination, error) {
	rows, err := db.Query(
		ctx,
		`SELECT id, environment_id, name, url, created_at, updated_at FROM destinations WHERE environment_id = $1`,
		environmentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var destinations []models.Destination

	for rows.Next() {
		var destination models.Destination

		err := rows.Scan(
			&destination.ID,
			&destination.EnvironmentID,
			&destination.Name,
			&destination.URL,
			&destination.CreatedAt,
			&destination.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, destination)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return destinations, nil
}

func UpdateDestinationByID(
	db *pgxpool.Pool,
	ctx context.Context,
	destinationID string,
	environmentID string,
	name string,
	url string,
) error {
	result, err := db.Exec(
		ctx,
		`UPDATE destinations SET name = $1, url = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND environment_id = $4`,
		name,
		url,
		destinationID,
		environmentID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func DeleteDestinationByID(
	db *pgxpool.Pool,
	ctx context.Context,
	destinationID string,
	environmentID string,
) error {
	result, err := db.Exec(
		ctx,
		`DELETE FROM destinations WHERE id = $1 AND environment_id = $2`,
		destinationID,
		environmentID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func GetDestinationForDelivery(
	db *pgxpool.Pool,
	ctx context.Context,
	destinationID string,
) (models.Destination, error) {
	var destination models.Destination
	err := db.QueryRow(
		ctx,
		`SELECT id, environment_id, name, url, created_at, updated_at FROM destinations WHERE id = $1`,
		destinationID,
	).Scan(
		&destination.ID,
		&destination.EnvironmentID,
		&destination.Name,
		&destination.URL,
		&destination.CreatedAt,
		&destination.UpdatedAt,
	)
	if err != nil {
		return models.Destination{}, err
	}

	return destination, nil
}
