package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IdempotencyRecord struct {
	ID             uuid.UUID  `db:"id"`
	IdempotencyKey string     `db:"idempotency_key"`
	UserID         uuid.UUID  `db:"user_id"`
	RequestPath    string     `db:"request_path"`
	RequestBody    string     `db:"request_body"`
	ResponseStatus int        `db:"response_status"`
	ResponseBody   string     `db:"response_body"`
	CreatedAt      time.Time  `db:"created_at"`
	ExpiresAt      time.Time  `db:"expires_at"`
}

func IdempotencyMiddleware(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only apply to POST, PUT, PATCH requests
		if c.Method() != "POST" && c.Method() != "PUT" && c.Method() != "PATCH" {
			return c.Next()
		}

		// Get idempotency key from header
		idempotencyKey := c.Get("Idempotency-Key")
		if idempotencyKey == "" {
			// Generate one based on request content for critical endpoints
			if isTransferEndpoint(c.Path()) {
				idempotencyKey = generateIdempotencyKey(c)
			} else {
				return c.Next()
			}
		}

		// Get user ID from context (set by auth middleware)
		userID, ok := c.Locals("userID").(uuid.UUID)
		if !ok {
			return c.Next()
		}

		// Check if request already processed
		var record IdempotencyRecord
		query := `SELECT * FROM idempotency_keys WHERE idempotency_key = $1 AND user_id = $2`
		err := db.Get(&record, query, idempotencyKey, userID)
		
		if err == nil {
			// Request already processed, return cached response
			if record.ExpiresAt.After(time.Now()) {
				c.Status(record.ResponseStatus)
				return c.SendString(record.ResponseBody)
			}
			// Expired, delete old record
			deleteQuery := `DELETE FROM idempotency_keys WHERE id = $1`
			db.Exec(deleteQuery, record.ID)
		}

		// Store request for processing
		requestBody := string(c.Body())
		insertQuery := `
			INSERT INTO idempotency_keys (idempotency_key, user_id, request_path, request_body, expires_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (idempotency_key) DO NOTHING`
		
		expiresAt := time.Now().Add(24 * time.Hour)
		_, err = db.Exec(insertQuery, idempotencyKey, userID, c.Path(), requestBody, expiresAt)
		if err != nil {
			// Concurrent request, return conflict
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Request is being processed. Please retry with a different idempotency key.",
			})
		}

		// Continue processing and capture response
		c.Locals("idempotencyKey", idempotencyKey)
		
		// Process the request
		err = c.Next()
		
		// Store the response
		responseBody := string(c.Response().Body())
		responseStatus := c.Response().StatusCode()
		
		updateQuery := `
			UPDATE idempotency_keys 
			SET response_status = $1, response_body = $2
			WHERE idempotency_key = $3 AND user_id = $4`
		
		db.Exec(updateQuery, responseStatus, responseBody, idempotencyKey, userID)
		
		return err
	}
}

func generateIdempotencyKey(c *fiber.Ctx) string {
	// Generate key based on user, path, and request body
	userID, _ := c.Locals("userID").(uuid.UUID)
	data := userID.String() + c.Path() + string(c.Body())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func isTransferEndpoint(path string) bool {
	return path == "/api/v1/transactions/transfer"
}