package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func ExecuteDelivery(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
) error {
	// Get the delivery and find its event and destination IDs
	delivery, err := repository.GetDeliveryForExecution(db, ctx, deliveryID)
	if err != nil {
		return err
	}

	// Get the event data that needs to be delivered
	event, err := repository.GetEventForDelivery(db, ctx, delivery.EventID)
	if err != nil {
		return err
	}

	// Get the destination url where the event data should be sent
	destination, err := repository.GetDestinationForDelivery(db, ctx, delivery.DestinationID)
	if err != nil {
		return err
	}

	// Convert the event data into JSON
	payload, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	// Send the event data to the destination url
	response, err := http.Post(
		destination.URL,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {

		deliveryError := err.Error()

		// Update delivery as failed because no HTTP response was recieved
		updateErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"failed",
			nil,
			nil,
			&deliveryError,
		)
		if updateErr != nil {
			return updateErr
		}

		return err
	}

	// Close the response body after reading it
	defer response.Body.Close()

	// Read response body returned by the destination
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Convert the response body from bytes to string
	responseBodystring := string(responseBody)

	if response.StatusCode >= 200 && response.StatusCode < 300 {

		// Update successful delivery result
		err := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"success",
			&response.StatusCode,
			&responseBodystring,
			nil,
		)
		if err != nil {
			return err
		}
		return nil
	}

	deliveryError := fmt.Sprintf("destination returned status %d", response.StatusCode)

	// Update delivery when the HTTP response is failed
	err = repository.UpdateDeliveryResults(
		db,
		ctx,
		delivery.ID,
		"failed",
		&response.StatusCode,
		&responseBodystring,
		&deliveryError,
	)
	if err != nil {
		return err
	}

	return fmt.Errorf("%s", deliveryError)
}
