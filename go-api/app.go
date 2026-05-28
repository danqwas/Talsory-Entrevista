package main

import (
	"bytes"
	"encoding/json"
	"go-api/config"
	"go-api/matrix"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

type RequestBody struct {
	Matrix [][]float64 `json:"matrix"`
}

type NodePayload struct {
	Rotated [][]float64 `json:"rotated"`
	MatrixQ [][]float64 `json:"matrix_q"`
	MatrixR [][]float64 `json:"matrix_r"`
}

func JWTMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No found Bearer token in Authorization header"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secretLimpio := strings.TrimSpace(config.JWT_SECRET)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretLimpio), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	return c.Next()
}

func SetupApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.FRONTEND_URL},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Authorization"},
	}))

	app.Post("/api/auth/login", func(c fiber.Ctx) error {
		secretLimpio := strings.TrimSpace(config.JWT_SECRET)

		claims := jwt.RegisteredClaims{
			Subject:   "interseguro_user",
			Issuer:    "go-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(secretLimpio))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		return c.JSON(fiber.Map{"token": tokenString})
	})

	app.Post("/api/matrix/process", JWTMiddleware, func(c fiber.Ctx) error {
		var body RequestBody

		if err := c.Bind().Body(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid matrix format",
			})
		}

		if len(body.Matrix) == 0 || len(body.Matrix[0]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "The matrix cannot be empty",
			})
		}

		matrixRotated := matrix.Rotate90Clockwise(body.Matrix)
		q, r := matrix.FactorizeQR(matrixRotated)

		payload := NodePayload{
			Rotated: matrixRotated,
			MatrixQ: q,
			MatrixR: r,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":        "Error preparing data for Node.js",
				"errorMessage": err.Error(),
				"data":         payload,
			})
		}

		req, err := http.NewRequest("POST", config.NODE_API_URL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":        "Error creating request to Node.js",
				"errorMessage": err.Error(),
				"data":         payload,
			})
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", c.Get("Authorization"))

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":        "Error communicating with Node.js",
				"errorMessage": err.Error(),
				"data":         payload,
			})
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		var nodeResponse interface{}
		_ = json.Unmarshal(respBody, &nodeResponse)

		return c.JSON(fiber.Map{
			"go_results": payload,
			"node_stats": nodeResponse,
		})
	})

	return app
}
