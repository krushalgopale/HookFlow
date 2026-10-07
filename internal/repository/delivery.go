package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
)

func SaveDelivery(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	eventID string,
	destinationID string,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO deliveries (id, event_id, destination_id, status) VALUES ($1, $2, $3, $4)`,
		deliveryID,
		eventID,
		destinationID,
		"pending",
	)

	return err
}

func GetDeliveryByEnvironment(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	environmentID string,
) (models.Delivery, error) {
	var delivery models.Delivery

	err := db.QueryRow(
		ctx,
		`SELECT d.id, d.event_id, d.destination_id, e.type, d.status, d.response_status, d.response_body, d.error, d.attempt_count, d.created_at, d.updated_at FROM deliveries d INNER JOIN events e ON d.event_id = e.id WHERE d.id = $1 AND e.environment_id = $2`,
		deliveryID,
		environmentID,
	).Scan(
		&delivery.ID,
		&delivery.EventID,
		&delivery.DestinationID,
		&delivery.EventType,
		&delivery.Status,
		&delivery.ResponseStatus,
		&delivery.ResponseBody,
		&delivery.Error,
		&delivery.AttemptCount,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)
	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}

func ListDeliveriesByEnvironment(
	db *pgxpool.Pool,
	ctx context.Context,
	environmentID string,
) ([]models.DeliveryListItem, error) {
	rows, err := db.Query(
		ctx,
		`SELECT d.status, e.type, d.created_at FROM deliveries d INNER JOIN events e ON d.event_id = e.id WHERE e.environment_id = $1 ORDER BY d.created_at DESC`,
		environmentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []models.DeliveryListItem

	for rows.Next() {
		var delivery models.DeliveryListItem

		err := rows.Scan(
			&delivery.Status,
			&delivery.EventType,
			&delivery.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, delivery)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func GetDeliveryByEvent(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	eventID string,
) (models.Delivery, error) {
	var delivery models.Delivery
	err := db.QueryRow(
		ctx,
		`SELECT d.status, e.type, d.created_at FROM deliveries d INNER JOIN events e ON d.event_id = e.id WHERE d.id = $1 AND e.event_id = $2 ORDER BY d.created_at DESC`,
		deliveryID,
		eventID,
	).Scan(
		&delivery.ID,
		&delivery.EventID,
		&delivery.DestinationID,
		&delivery.EventType,
		&delivery.Status,
		&delivery.ResponseStatus,
		&delivery.ResponseBody,
		&delivery.Error,
		&delivery.AttemptCount,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)
	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}

func ListDeliveriesByEvent(
	db *pgxpool.Pool,
	ctx context.Context,
	eventID string,
) ([]models.DeliveryListItem, error) {
	rows, err := db.Query(
		ctx,
		`SELECT  d.status, e.type, d.created_at FROM deliveries d INNER JOIN events e ON d.event_id = e.id WHERE d.event_id = $1 ORDER BY d.created_at DESC`,
		eventID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var deliveries []models.DeliveryListItem
	for rows.Next() {
		var delivery models.DeliveryListItem

		err := rows.Scan(
			&delivery.Status,
			&delivery.EventType,
			&delivery.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, delivery)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func UpdateDeliveryResults(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	status string,
	responseStatus *int,
	responseBody *string,
	deliveryError *string,
) error {
	result, err := db.Exec(
		ctx,
		`UPDATE deliveries SET status = $1, response_status = $2, response_body = $3, error = $4, attempt_count = attempt_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $5`,
		status,
		responseStatus,
		responseBody,
		deliveryError,
		deliveryID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func GetDeliveryForExecution(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
) (models.Delivery, error) {
	var delivery models.Delivery
	err := db.QueryRow(
		ctx,
		`SELECT id, event_id, destination_id, status FROM deliveries WHERE id = $1`,
		deliveryID,
	).Scan(
		&delivery.ID,
		&delivery.EventID,
		&delivery.DestinationID,
		&delivery.Status,
	)
	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}
