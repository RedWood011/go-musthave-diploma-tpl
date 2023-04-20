package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	transport "github.com/RedWood011/cmd/gophermart/internal/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)

func TestUserRegistration(t *testing.T) {
	existUser := "userExist"
	password := "0123456789"
	app, err := initTest()
	require.NoError(t, err)
	testTable := []struct {
		name       string
		login      string
		password   string
		statusCode int
	}{
		{
			name:       "Error registration user",
			login:      "test",
			password:   "test",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Registration successful",
			login:      "userSuccess",
			password:   "0123456789",
			statusCode: http.StatusOK,
		},
		{
			name:       "User exists",
			login:      existUser,
			password:   password,
			statusCode: http.StatusConflict,
		},
	}

	resp := createUser(t, app, existUser, password)
	defer resp.Body.Close()
	require.Equal(t, resp.StatusCode, http.StatusOK)

	for _, testCases := range testTable {
		t.Run(testCases.name, func(t *testing.T) {
			response := createUser(t, app, testCases.login, testCases.password)
			assert.Equal(t, response.StatusCode, testCases.statusCode)
			response.Body.Close()
		})
	}
}
func TestUserAuthorization(t *testing.T) {
	app, err := initTest()
	require.NoError(t, err)
	User := "AuthUser"
	password := "0123456789"

	testTable := []struct {
		name       string
		login      string
		password   string
		statusCode int
	}{
		{
			name:       "Error authorization user:invalid login or password",
			login:      "test",
			password:   "test",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "Authorization user success",
			login:      User,
			password:   password,
			statusCode: http.StatusOK,
		},
	}

	resp := createUser(t, app, User, password)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	for _, testCases := range testTable {
		t.Run(testCases.name, func(t *testing.T) {
			response := authUser(t, app, testCases.login, testCases.password)
			assert.Equal(t, response.StatusCode, testCases.statusCode)
		})
	}

}

func createUser(t *testing.T, app *fiber.App, login string, password string) *http.Response {
	expected := transport.UserRegRequest{
		Login:    login,
		Password: password,
	}
	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api/user/register", createReqBody(t, expected))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 1000)
	require.NoError(t, err)
	return resp
}
func authUser(t *testing.T, app *fiber.App, login string, password string) *http.Response {
	expected := transport.UserRegRequest{
		Login:    login,
		Password: password,
	}
	req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api/user/login", createReqBody(t, expected))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 1000)
	require.NoError(t, err)
	return resp

}
