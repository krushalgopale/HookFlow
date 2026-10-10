package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
)

// Deliveries
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

func DeliveryBelongsToUser(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	userID string,
) (bool, error) {
	var belongs bool

	err := db.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM deliveries d INNER JOIN events e ON e.id = d.event_id INNER JOIN environments env ON env.id = e.environment_id INNER JOIN tenants t ON t.id = env.tenant_id WHERE d.id = $1 AND t.user_id = $2)`,
		deliveryID,
		userID,
	).Scan(&belongs)
	if err != nil {
		return false, err
	}

	return belongs, nil
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
		`SELECT d.id, d.event_id, d.destination_id, e.type, dest.name, d.status, d.response_status, d.response_body, d.error, d.attempt_count, d.created_at, d.updated_at FROM deliveries d INNER JOIN destinations dest ON d.destination_id = dest.id INNER JOIN events e ON d.event_id = e.id WHERE d.id = $1 AND e.environment_id = $2`,
		deliveryID,
		environmentID,
	).Scan(
		&delivery.ID,
		&delivery.EventID,
		&delivery.DestinationID,
		&delivery.EventType,
		&delivery.DestinationName,
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
	limit int,
	offset int,
) ([]models.DeliveryListItem, error) {
	rows, err := db.Query(
		ctx,
		`SELECT d.status, e.type, dest.name, d.created_at FROM deliveries d INNER JOIN destinations dest ON dest.id = d.destination_id INNER JOIN events e ON d.event_id = e.id WHERE e.environment_id = $1 ORDER BY d.created_at DESC LIMIT $2 OFFSET $3`,
		environmentID,
		limit,
		offset,
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
			&delivery.DestinationName,
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
		`SELECT d.id, d.event_id, d.destination_id, e.type, dest.name, d.status, d.response_status, d.response_body, d.error, d.attempt_count, d.created_at, d.updated_at FROM deliveries d INNER JOIN destinations dest ON d.destination_id = dest.id INNER JOIN events e ON d.event_id = e.id WHERE d.id = $1 AND d.event_id = $2`,
		deliveryID,
		eventID,
	).Scan(
		&delivery.ID,
		&delivery.EventID,
		&delivery.DestinationID,
		&delivery.EventType,
		&delivery.DestinationName,
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
		`SELECT  d.status, e.type, dest.name, d.created_at FROM deliveries d INNER JOIN destinations dest ON d.destination_id = dest.id INNER JOIN events e ON d.event_id = e.id WHERE d.event_id = $1 ORDER BY d.created_at DESC`,
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
			&delivery.DestinationName,
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

// Delivery Execution
func UpdateDeliveryResults(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	status string,
	responseStatus *int,
	responseBody *string,
	deliveryError *string,
	nextAttemptAt *time.Time,
) error {
	result, err := db.Exec(
		ctx,
		`UPDATE deliveries SET status = $1, response_status = $2, response_body = $3, error = $4, attempt_count = attempt_count + 1, next_attempt_at = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6`,
		status,
		responseStatus,
		responseBody,
		deliveryError,
		nextAttemptAt,
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
		`SELECT id, event_id, destination_id, status, max_attempts, next_attempt_at FROM deliveries WHERE id = $1`,
		deliveryID,
	).Scan(
		&delivery.ID,
		&delivery.EventID,
		&delivery.DestinationID,
		&delivery.Status,
		&delivery.MaxAttempts,
		&delivery.NextAttemptAt,
	)
	if err != nil {
		return models.Delivery{}, err
	}

	return delivery, nil
}

// Delivery Attempts

func SaveDeliveryAttempt(
	db *pgxpool.Pool,
	ctx context.Context,
	attempt models.DeliveryAttempt,
) error {
	_, err := db.Exec(
		ctx,
		`INSERT INTO delivery_attempts (id, delivery_id, attempt_number, status, response_status, response_body, error) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		attempt.ID,
		attempt.DeliveryID,
		attempt.AttemptNumber,
		attempt.Status,
		attempt.ResponseStatus,
		attempt.ResponseBody,
		attempt.Error,
	)
	if err != nil {
		return err
	}

	return nil
}

func ListDeliveryAttemptByDeliveryID(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	limit int,
	offset int,
) ([]models.DeliveryAttemptListItem, error) {
	rows, err := db.Query(
		ctx,
		`SELECT dest.name, da.attempt_number, da.status, da.created_at FROM delivery_attempts da INNER JOIN deliveries d ON d.id = da.delivery_id INNER JOIN destinations dest ON d.destination_id = dest.id WHERE da.delivery_id = $1 ORDER BY da.attempt_number ASC LIMIT $2 OFFSET $3`,
		deliveryID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveryAttempts []models.DeliveryAttemptListItem

	for rows.Next() {
		var deliveryAttempt models.DeliveryAttemptListItem

		err := rows.Scan(
			&deliveryAttempt.DestinationName,
			&deliveryAttempt.AttemptNumber,
			&deliveryAttempt.Status,
			&deliveryAttempt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		deliveryAttempts = append(deliveryAttempts, deliveryAttempt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deliveryAttempts, nil
}
