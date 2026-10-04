package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
)

func SaveEvent(
	db *pgxpool.Pool,
	ctx context.Context,
	eventID string,
	environmentID string,
	eventType string,
	data []byte,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO events (id, environment_id, type, data) VALUES ($1, $2, $3, $4)`,
		eventID,
		environmentID,
		eventType,
		data,
	)
	return err
}

func GetEvent(
	db *pgxpool.Pool,
	ctx context.Context,
	eventID string,
	environmentID string,
) (*models.Event, error) {
	var event models.Event

	err := db.QueryRow(
		ctx,
		`SELECT id, environment_id, type, data, created_at FROM events WHERE id = $1 AND environment_id = $2`,
		eventID,
		environmentID,
	).Scan(
		&event.ID, &event.EnvironmentID, &event.Type, &event.Data, &event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func ListEvents(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
	limit int,
	offset int,
) ([]models.Event, int, error) {
	events := []models.Event{}

	// Database Operation with Pagination
	rows, err := db.Query(
		ctx,
		`SELECT id, environment_id, type, data, created_at FROM events WHERE environment_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		environmentID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID, &event.EnvironmentID, &event.Type, &event.Data, &event.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		events = append(events, event)
	}

	var total int

	err = db.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM events`,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return events, total, nil
}
