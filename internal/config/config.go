// Package config provides functionality for loading and managing application configuration.
// It retrieves configuration values from environment variables and ensures that all required
// settings are properly initialized. This package is essential for setting up the application's
// runtime environment, including database connections, JWT secrets, and server settings.
package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config represents the application configuration.
// It contains environment-specific settings such as the environment name,
// server port, JWT secret, and database URL.
type Config struct {
	Env                string // The current environment (e.g., "dev", "prod").
	Port               string // The port on which the server will run.
	FrontendURL        string // The URL of the frontend application.
	AuthServiceURL     string // The URL of the authentication service.
	TemplateServiceURL string // The URL of the dashboard service.
	PDFServiceURL      string // The URL of the PDF service.
	JWTSecret          []byte // The secret key used for signing JWT tokens.
	CookieSecure       bool   // The secure flag for cookies (true for HTTPS, false for HTTP).
}

const (
	envKey             = "ENV"                  // Environment variable key for the environment name.
	portEnv            = "PORT"                 // Environment variable key for the server port.
	frontEndKey        = "FRONTEND_URL"         // Environment variable key for the frontend URL.
	authServiceKey     = "AUTH_SERVICE_URL"     // Environment variable key for the authentication service URL.
	templateServiceKey = "TEMPLATE_SERVICE_URL" // Environment variable key for the dashboard service URL.
	pdfServiceKey      = "PDF_SERVICE_URL"      // Environment variable key for the PDF service URL.
	jwtSecretKey       = "JWT_SECRET"           // Environment variable key for the JWT secret.
	cookieSecureKey    = "COOKIE_SECURE"        // Environment variable key for the secure flag of cookies.

	defaultEnvKey = "dev" // Default environment name if none is provided.
)

// Load retrieves the application configuration from environment variables.
// It ensures that all required configuration values are set and returns an error
// if any mandatory value is missing.
//
// Returns:
//   - Config: The loaded application configuration.
//   - error: An error if any required configuration value is missing.
func Load() (Config, error) {
	var c Config

	c.Env = os.Getenv(envKey)
	if c.Env == "" {
		c.Env = defaultEnvKey
	}

	c.Port = getEnv(portEnv)
	if c.Port == "" {
		return Config{}, errors.New("empty key")
	}

	c.FrontendURL = getEnv(frontEndKey)
	if c.FrontendURL == "" {
		return Config{}, errors.New("empty key")
	}

	c.AuthServiceURL = getEnv(authServiceKey)
	if c.AuthServiceURL == "" {
		return Config{}, errors.New("empty key")
	}

	c.TemplateServiceURL = getEnv(templateServiceKey)
	if c.TemplateServiceURL == "" {
		return Config{}, errors.New("empty key")
	}

	c.PDFServiceURL = getEnv(pdfServiceKey)
	if c.PDFServiceURL == "" {
		return Config{}, errors.New("empty key")
	}

	c.JWTSecret = []byte(getEnv(jwtSecretKey))
	if len(c.JWTSecret) == 0 {
		return Config{}, errors.New("empty key")
	}

	var err error
	c.CookieSecure, err = strconv.ParseBool(getEnv(cookieSecureKey))
	if err != nil {
		return Config{}, errors.New("invalid value for COOKIE_SECURE:" + err.Error())
	}

	return c, nil
}

// getEnv retrieves the value of an environment variable.
// If the variable is not set, it logs an error and returns an empty string.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//
// Returns:
//   - string: The value of the environment variable, or an empty string if not set.
func getEnv(key string) string {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	val := os.Getenv(key)
	if val == "" {
		log.Error().Str("var", key).Msg("Failed to load environment")
	}
	return val
}
