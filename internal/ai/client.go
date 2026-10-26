package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/krushalgopale/HookFlow/internal/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Analyze(ctx context.Context, request models.FailureRequest) (models.FailureResponse, error) {
	var result models.FailureResponse

	// JSON marshalling
	body, err := json.Marshal(request)
	if err != nil {
		return result, fmt.Errorf("Failed to marshal AI request: %w", err)
	}

	// Post request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/analyze",
		bytes.NewReader(body),
	)
	if err != nil {
		return result, fmt.Errorf("Failed to create AI request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("Failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf(
			"AI service returned HTTP %d",
			resp.StatusCode,
		)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("Failed to decode AI response: %w", err)
	}

	return result, nil
}
