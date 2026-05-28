package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	originalPort := os.Getenv("PORT")
	originalUrl := os.Getenv("NODE_API_URL")
	originalSecret := os.Getenv("JWT_SECRET")
	originalFrontend := os.Getenv("FRONTEND_URL")

	defer func() {

		os.Setenv("PORT", originalPort)
		os.Setenv("NODE_API_URL", originalUrl)
		os.Setenv("JWT_SECRET", originalSecret)
		os.Setenv("FRONTEND_URL", originalFrontend)
	}()

	t.Run("Debe cargar valores por defecto si no hay variables de entorno", func(t *testing.T) {

		os.Unsetenv("PORT")
		os.Unsetenv("NODE_API_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("FRONTEND_URL")

		LoadConfig()

		if PORT != "3000" {
			t.Errorf("Se esperaba el puerto 3000 por defecto, se obtuvo %s", PORT)
		}
		if NODE_API_URL != "http://localhost:4000/api/statistics" {
			t.Errorf("Se esperaba la URL por defecto, se obtuvo %s", NODE_API_URL)
		}
		if JWT_SECRET != "MiClaveSecretaSuperSegura123!" {
			t.Errorf("Se esperaba el secreto por defecto, se obtuvo %s", JWT_SECRET)
		}
		if FRONTEND_URL != "http://localhost:5173" {
			t.Errorf("Se esperaba la URL del frontend por defecto, se obtuvo %s", FRONTEND_URL)
		}
	})

	t.Run("Debe leer las variables de entorno si estan presentes", func(t *testing.T) {

		os.Setenv("PORT", "8080")
		os.Setenv("NODE_API_URL", "http://mi-api.com/stats")
		os.Setenv("JWT_SECRET", "UltraSecreto")
		os.Setenv("FRONTEND_URL", "http://localhost:5173")

		LoadConfig()

		if PORT != "8080" {
			t.Errorf("Se esperaba el puerto 8080, se obtuvo %s", PORT)
		}
		if NODE_API_URL != "http://mi-api.com/stats" {
			t.Errorf("Se esperaba la URL inyectada, se obtuvo %s", NODE_API_URL)
		}
		if JWT_SECRET != "UltraSecreto" {
			t.Errorf("Se esperaba el secreto inyectado, se obtuvo %s", JWT_SECRET)
		}
		if FRONTEND_URL != "http://localhost:5173" {
			t.Errorf("Se esperaba la URL del frontend inyectada, se obtuvo %s", FRONTEND_URL)
		}
	})
}
