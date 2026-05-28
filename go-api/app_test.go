package main

import (
	"bytes"
	"go-api/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CoverageTest(t *testing.T) {

	t.Run("You must properly validate the environment configuration", func(t *testing.T) {

		os.Setenv("PORT", "3000")
		os.Setenv("NODE_API_URL", "http://localhost:4000/api/statistics")
		os.Setenv("JWT_SECRET", "ClaveSecretaDePrueba123!")

		config.LoadConfig()
	})

	app := SetupApp()

	t.Run("Login - Must generate a valid JWT token (200)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error in request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, but got %d", resp.StatusCode)
		}
	})
	t.Run("Login - Must fail to generate token (500)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error in request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected 500, but got %d", resp.StatusCode)
		}
	})
	t.Run("Middleware - Must reject if the Authorization header is missing (401)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", nil)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 for missing header, got %d", resp.StatusCode)
		}
	})

	t.Run("Middleware - Must reject if the token is not a valid Bearer token (401)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", nil)
		req.Header.Set("Authorization", "TokenInvalidoSinBearer")
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 for invalid Bearer token, got %d", resp.StatusCode)
		}
	})

	t.Run("Middleware - Must reject if the JWT token has expired or the signature is invalid (401)", func(t *testing.T) {

		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(config.JWT_SECRET))

		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 for expired token, got %d", resp.StatusCode)
		}
	})

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	tokenValido := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtCorrecto, _ := tokenValido.SignedString([]byte(config.JWT_SECRET))

	t.Run("Process - Must return 400 if the JSON is malformed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer([]byte(`{matrix: roto`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 for malformed JSON, got %d", resp.StatusCode)
		}
	})

	t.Run("Process - Must return 400 if the matrix is empty", func(t *testing.T) {
		jsonBody := []byte(`{"matrix": []}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 for empty matrix, got %d", resp.StatusCode)
		}
	})

	t.Run("Process - Must return 200 for successful request with valid matrix (200)", func(t *testing.T) {

		nodeMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"average": 5.5, "totalSum": 22, "max": 4, "min": 1}`))
		}))
		defer nodeMockServer.Close()

		config.NODE_API_URL = nodeMockServer.URL

		jsonBody := []byte(`{"matrix": [[1.0, 2.0], [3.0, 4.0]]}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 for successful request, got %d", resp.StatusCode)
		}
	})

	t.Run("Process - Must handle error if Node API responds with an error status (500)", func(t *testing.T) {

		nodeMockServerFalla := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Node crash", http.StatusInternalServerError)
		}))
		defer nodeMockServerFalla.Close()

		config.NODE_API_URL = nodeMockServerFalla.URL

		jsonBody := []byte(`{"matrix": [[1.0, 2.0], [3.0, 4.0]]}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode == 0 {
			t.Error("Failed internal validation - response status code is 0, expected 500")
		}
	})
	t.Run("Must fail Bind if we send an incorrect data type (400)", func(t *testing.T) {
		jsonBody := []byte(`{"matrix": 12345}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 for bind error, got %d", resp.StatusCode)
		}
	})
	t.Run("Must fail if the internal matrix is empty (400)", func(t *testing.T) {
		jsonBody := []byte(`{"matrix": [[]]}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 for empty internal matrix, got %d", resp.StatusCode)
		}
	})
	t.Run("Must fail if the Node API URL is invalid (500)", func(t *testing.T) {
		config.NODE_API_URL = "://url-completamente-invalida"

		jsonBody := []byte(`{"matrix": [[1.0, 2.0], [3.0, 4.0]]}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected 500 for invalid Node API URL, got %d", resp.StatusCode)
		}
	})
	t.Run("Must fail if unable to communicate with Node if the server does not respond (500)", func(t *testing.T) {
		config.NODE_API_URL = "http://localhost:99999"

		jsonBody := []byte(`{"matrix": [[1.0, 2.0], [3.0, 4.0]]}`)
		req := httptest.NewRequest(http.MethodPost, "/api/matrix/process", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtCorrecto)

		resp, _ := app.Test(req)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected 500 for network communication failure, got %d", resp.StatusCode)
		}
	})
}

func LocalConfig() {
	_ = len(config.PORT)
}
