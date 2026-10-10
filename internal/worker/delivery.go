package worker

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

const deliveryTimeout = 10 * time.Second

func ExecuteDelivery(
	db *pgxpool.Pool,
	ctx context.Context,
	deliveryID string,
	attempt int,
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

	mac := hmac.New(sha256.New, []byte(destination.SigningSecret))
	mac.Write(payload)

	signature := hex.EncodeToString(mac.Sum(nil))

	// Set response waiting time
	client := &http.Client{
		Timeout: deliveryTimeout,
	}

	// Send the event data to the destination url
	request, err := http.NewRequest(
		http.MethodPost,
		destination.URL,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "HookFlow/1.0")
	request.Header.Set("X-HookFlow-Delivery-ID", delivery.ID)
	request.Header.Set("X-HookFlow-Event-ID", delivery.EventID)
	request.Header.Set("X-HookFlow-Environment-ID", event.EnvironmentID)
	request.Header.Set("X-HookFlow-Signature", "sha256="+signature)

	response, err := client.Do(request)
	if err != nil {
		requestErr := err
		deliveryError := err.Error()

		// Save each delivery attempts
		attemptRecord := models.DeliveryAttempt{
			ID:             "att_" + uuid.New().String(),
			DeliveryID:     delivery.ID,
			AttemptNumber:  attempt,
			Status:         "failed",
			ResponseStatus: nil,
			ResponseBody:   nil,
			Error:          &deliveryError,
		}

		SaveErr := repository.SaveDeliveryAttempt(db, ctx, attemptRecord)
		if SaveErr != nil {
			return SaveErr
		}

		var nextAttemptAt *time.Time

		// Store the time when the next retry should happen
		if attempt < delivery.MaxAttempts {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			jitter := time.Duration(rand.Float64() * float64(time.Second))
			next := time.Now().Add(delay + jitter)
			nextAttemptAt = &next
		}

		// Update delivery as failed because no HTTP response was recieved
		updateErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"failed",
			nil,
			nil,
			&deliveryError,
			nextAttemptAt,
		)
		if updateErr != nil {
			return updateErr
		}

		// Retry if more attempts are available
		if attempt < delivery.MaxAttempts {
			if nextAttemptAt != nil {
				time.Sleep(time.Until(*nextAttemptAt))
			}

			return ExecuteDelivery(
				db,
				ctx,
				delivery.ID,
				attempt+1,
			)
		}

		return requestErr
	}

	// Close the response body after reading it
	defer response.Body.Close()

	// Read response body returned by the destination
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		deliveryError := fmt.Sprintf("failed to read destination response body: %v", err)

		attemptRecord := models.DeliveryAttempt{
			ID:             "att_" + uuid.New().String(),
			DeliveryID:     delivery.ID,
			AttemptNumber:  attempt,
			Status:         "failed",
			ResponseStatus: &response.StatusCode,
			Error:          &deliveryError,
		}

		saveErr := repository.SaveDeliveryAttempt(db, ctx, attemptRecord)
		if saveErr != nil {
			return saveErr
		}

		var nextAttemptAt *time.Time
		if attempt < delivery.MaxAttempts {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			jitter := time.Duration(rand.Float64() * float64(time.Second))
			next := time.Now().Add(delay + jitter)
			nextAttemptAt = &next
		}

		updateErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"failed",
			&response.StatusCode,
			nil,
			&deliveryError,
			nextAttemptAt,
		)
		if updateErr != nil {
			return updateErr
		}

		if attempt < delivery.MaxAttempts {
			if nextAttemptAt != nil {
				time.Sleep(time.Until(*nextAttemptAt))
			}

			return ExecuteDelivery(db, ctx, delivery.ID, attempt+1)
		}

		return fmt.Errorf("%s", deliveryError)

	}

	// Convert the response body from bytes to string
	responseBodyString := string(responseBody)

	// 2xx Sucess
	if response.StatusCode >= 200 && response.StatusCode < 300 {

		// Save each delivery attempts
		attemptRecord := models.DeliveryAttempt{
			ID:             "att_" + uuid.New().String(),
			DeliveryID:     delivery.ID,
			AttemptNumber:  attempt,
			Status:         "success",
			ResponseStatus: &response.StatusCode,
			ResponseBody:   &responseBodyString,
			Error:          nil,
		}

		saveErr := repository.SaveDeliveryAttempt(db, ctx, attemptRecord)
		if saveErr != nil {
			return saveErr
		}

		// Update successful delivery result
		updateErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"success",
			&response.StatusCode,
			&responseBodyString,
			nil,
			nil,
		)
		if updateErr != nil {
			return updateErr
		}
		return nil
	}

	// 4xx Client Error
	if response.StatusCode >= 400 && response.StatusCode < 500 &&
		response.StatusCode != http.StatusTooManyRequests {

		deliveryError := fmt.Sprintf(
			"destination returned status %d",
			response.StatusCode,
		)

		// Save each delivery attempts
		attemptRecord := models.DeliveryAttempt{
			ID:             "att_" + uuid.New().String(),
			DeliveryID:     delivery.ID,
			AttemptNumber:  attempt,
			Status:         "failed",
			ResponseStatus: &response.StatusCode,
			ResponseBody:   &responseBodyString,
			Error:          &deliveryError,
		}

		saveErr := repository.SaveDeliveryAttempt(db, ctx, attemptRecord)
		if saveErr != nil {
			return saveErr
		}

		UpdateErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"failed",
			&response.StatusCode,
			&responseBodyString,
			&deliveryError,
			nil,
		)
		if UpdateErr != nil {
			return UpdateErr
		}

		return fmt.Errorf("%s", deliveryError)
	}

	// 429 too many requests
	if response.StatusCode == http.StatusTooManyRequests {
		deliveryError := fmt.Sprintf("destination returned status %d", response.StatusCode)

		// Save each delivery attempts
		attemptRecord := models.DeliveryAttempt{
			ID:             "att_" + uuid.New().String(),
			DeliveryID:     delivery.ID,
			AttemptNumber:  attempt,
			Status:         "failed",
			ResponseStatus: &response.StatusCode,
			ResponseBody:   &responseBodyString,
			Error:          &deliveryError,
		}

		saveErr := repository.SaveDeliveryAttempt(db, ctx, attemptRecord)
		if saveErr != nil {
			return saveErr
		}

		var nextAttemptAt *time.Time

		// Store the time when the next retry should happen
		if attempt < delivery.MaxAttempts {
			retryAfter := response.Header.Get("Retry-After")

			// Use the server's Retry-After value when provided
			if retryAfter != "" {
				// for seconds
				seconds, err := strconv.Atoi(retryAfter)
				if err == nil && seconds >= 0 {
					next := time.Now().Add(time.Duration(seconds) * time.Second)
					nextAttemptAt = &next
				}

				// for date and time
				if nextAttemptAt == nil {
					retryTime, err := http.ParseTime(retryAfter)
					if err == nil {
						nextAttemptAt = &retryTime
					}
				}
			}

			// Use exponential backoff with jitter when Retry-After is not provided
			if nextAttemptAt == nil {
				delay := time.Duration(1<<uint(attempt-1)) * time.Second
				jitter := time.Duration(rand.Float64() * float64(time.Second))
				next := time.Now().Add(delay + jitter)
				nextAttemptAt = &next
			}
		}

		// Update delivery when the destination HTTP response is failed
		updtaeErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"failed",
			&response.StatusCode,
			&responseBodyString,
			&deliveryError,
			nextAttemptAt,
		)
		if updtaeErr != nil {
			return updtaeErr
		}

		// Retry if more attempts are available
		if attempt < delivery.MaxAttempts {
			if nextAttemptAt != nil {
				time.Sleep(time.Until(*nextAttemptAt))
			}
			return ExecuteDelivery(
				db,
				ctx,
				delivery.ID,
				attempt+1,
			)
		}

		return fmt.Errorf("%s", deliveryError)
	}

	// 5xx Server Error
	if response.StatusCode >= 500 && response.StatusCode < 600 {
		deliveryError := fmt.Sprintf("destination returned status %d", response.StatusCode)

		// Save each delivery attempts
		attemptRecord := models.DeliveryAttempt{
			ID:             "att_" + uuid.New().String(),
			DeliveryID:     delivery.ID,
			AttemptNumber:  attempt,
			Status:         "failed",
			ResponseStatus: &response.StatusCode,
			ResponseBody:   &responseBodyString,
			Error:          &deliveryError,
		}

		saveErr := repository.SaveDeliveryAttempt(db, ctx, attemptRecord)
		if saveErr != nil {
			return saveErr
		}

		var nextAttemptAt *time.Time

		// Store the time when the next retry should happen
		if attempt < delivery.MaxAttempts {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			jitter := time.Duration(rand.Float64() * float64(time.Second))
			next := time.Now().Add(delay + jitter)
			nextAttemptAt = &next
		}

		// Update delivery when the destination HTTP response is failed
		updateErr := repository.UpdateDeliveryResults(
			db,
			ctx,
			delivery.ID,
			"failed",
			&response.StatusCode,
			&responseBodyString,
			&deliveryError,
			nextAttemptAt,
		)
		if updateErr != nil {
			return updateErr
		}

		// Retry if more attempts are available
		if attempt < delivery.MaxAttempts {
			if nextAttemptAt != nil {
				time.Sleep(time.Until(*nextAttemptAt))
			}
			return ExecuteDelivery(
				db,
				ctx,
				delivery.ID,
				attempt+1,
			)
		}

		return fmt.Errorf("%s", deliveryError)
	}
	return nil
}
