package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedWood011/cmd/gophermart/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareJWT(t *testing.T) {
	app := fiber.New()

	cfg := config.TokenConfig{
		SecretKey: "mysecretkey",
		UserKey:   "user_id",
	}

	app.Use(MiddlewareJWT(cfg))

	app.Get("/api/protected", func(c *fiber.Ctx) error {
		userID := c.Locals(cfg.UserKey).(string)
		return c.JSON(fiber.Map{
			"message": "Hello, " + userID,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	resp, err := app.Test(req)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		cfg.UserKey: "123",
	})
	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		t.Fatalf("Error creating test token: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err = app.Test(req)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	var data fiber.Map
	err = json.NewDecoder(resp.Body).Decode(&data)
	require.NoError(t, err)
	if message := data["message"].(string); message != "Hello, 123" {
		t.Errorf("Expected message 'Hello, 123' but got '%s'", message)
	}
}
