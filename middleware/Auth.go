package middleware

import (
	"crypto/sha256"
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

type AuthMiddleware struct {
	hashedAPIKey [32]byte
	Middleware   func(*fiber.Ctx) error
}

func NewAuthMiddleware(apiKey string) *AuthMiddleware {
	authMiddlware := &AuthMiddleware{hashedAPIKey: sha256.Sum256([]byte(apiKey))}
	authMiddlware.Middleware = keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			hashedKey := sha256.Sum256([]byte(key))

			if subtle.ConstantTimeCompare(hashedKey[:], authMiddlware.hashedAPIKey[:]) == 1 {
				return true, nil
			}
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
				"data":    nil,
			})
		},
	})
	return authMiddlware
}
