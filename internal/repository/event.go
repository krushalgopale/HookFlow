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
	eventType string,
	data []byte,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO events (id, type, data) VALUES ($1, $2, $3)`,
		eventID,
		eventType,
		data,
	)
	return err
}

func GetEvent(
	db *pgxpool.Pool,
	ctx context.Context,
	eventID string,
) (*models.Event, error) {
	var event models.Event

	err := db.QueryRow(
		ctx,
		`SELECT id, type, data, created_at FROM events WHERE id = $1`,
		eventID,
	).Scan(
		&event.ID, &event.Type, &event.Data, &event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func ListEvents(
	db *pgxpool.Pool,
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Event,int, error) {
	events := []models.Event{}

	// Database Operation with Pagination
	rows, err := db.Query(
		ctx,
		`SELECT id, type, data, created_at FROM events ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
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
			&event.ID, &event.Type, &event.Data, &event.CreatedAt,
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
