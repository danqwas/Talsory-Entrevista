package config

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var (
	PORT         string
	NODE_API_URL string
	JWT_SECRET   string
	FRONTEND_URL string
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Environment variables loaded from system environment (no .env file found)")
	}

	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "3000"
	}

	envNodeURL := os.Getenv("NODE_API_URL")
	if envNodeURL == "" {
		envNodeURL = "http://localhost:4000/api/statistics"
	}

	envJWTSecret := os.Getenv("JWT_SECRET")
	if envJWTSecret == "" {
		envJWTSecret = "MiClaveSecretaSuperSegura123!"
	}

	envFrontendURL := os.Getenv("FRONTEND_URL")
	if envFrontendURL == "" {
		envFrontendURL = "http://localhost:5173"
	}

	type ConfigValidator struct {
		PORT         string `validate:"required,numeric,gte=1,lte=65535"`
		NODE_API_URL string `validate:"required,url"`
		JWT_SECRET   string `validate:"required"`
		FRONTEND_URL string `validate:"required,url"`
	}

	cfgData := ConfigValidator{
		PORT:         envPort,
		NODE_API_URL: envNodeURL,
		JWT_SECRET:   envJWTSecret,
		FRONTEND_URL: envFrontendURL,
	}

	validate := validator.New()
	if err := validate.Struct(cfgData); err != nil {
		slog.Error("Critical validation error", "error", err)
		os.Exit(1)
	}

	PORT = cfgData.PORT
	NODE_API_URL = cfgData.NODE_API_URL
	JWT_SECRET = cfgData.JWT_SECRET
	FRONTEND_URL = cfgData.FRONTEND_URL
	slog.Info("Configuration loaded successfully", "PORT", PORT, "NODE_API_URL", NODE_API_URL, "FRONTEND_URL", FRONTEND_URL)
}
