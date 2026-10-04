package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func VerifyAPIKey(
	db *pgxpool.Pool,
	ctx context.Context,
	key string,
) (*models.APIKey, error) {
	// Hash the api key received from the client
	hash := sha256.Sum256([]byte(key))
	// Convert the sha256 hash into a 64 character
	keyHash := hex.EncodeToString(hash[:])

	// Look up the api key using the hashed value stored in database
	apiKey, err := repository.GetAPIKeyByKey(db, ctx, keyHash)
	if err != nil {
		return nil, err
	}

	// check whether key has expired or not
	if apiKey.ExpiresAt != nil && !apiKey.ExpiresAt.After(time.Now()) {
		return nil, fmt.Errorf("api key has expired")
	}

	return apiKey, nil
}
