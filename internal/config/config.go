// Package config provides functionality for loading and managing application configuration.
// It retrieves configuration values from environment variables and ensures that all required
// settings are properly initialized. This package is essential for setting up the application's
// runtime environment, including database connections, JWT secrets, and server settings.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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
	// Initialize logger once, if it's specific to config loading or used by getEnv.
	// However, it's better if getEnv doesn't log directly but returns errors or relies on a global logger.
	// For now, let's assume the global logger is initialized in main.
	// If getEnv needs to log, it should use the global logger.
	// The per-call logger setup in getEnv was problematic.

	var c Config

	c.Env = os.Getenv(envKey)
	if c.Env == "" {
		c.Env = defaultEnvKey
	}

	c.Port = getEnv(portEnv, true)
	if c.Port == "" {
		return Config{}, errors.New("empty key: " + portEnv)
	}

	c.FrontendURL = getEnv(frontEndKey, true)
	if c.FrontendURL == "" {
		return Config{}, errors.New("empty key: " + frontEndKey)
	}

	c.AuthServiceURL = getEnv(authServiceKey, true)
	if c.AuthServiceURL == "" {
		return Config{}, errors.New("empty key: " + authServiceKey)
	}

	c.TemplateServiceURL = getEnv(templateServiceKey, true)
	if c.TemplateServiceURL == "" {
		return Config{}, errors.New("empty key: " + templateServiceKey)
	}

	c.PDFServiceURL = getEnv(pdfServiceKey, true)
	if c.PDFServiceURL == "" {
		return Config{}, errors.New("empty key: " + pdfServiceKey)
	}

	c.JWTSecret = []byte(getEnv(jwtSecretKey, true))
	if len(c.JWTSecret) == 0 {
		return Config{}, errors.New("empty key: " + jwtSecretKey)
	}

	var err error
	cookieSecureStr := getEnv(cookieSecureKey, true)
	if cookieSecureStr == "" { // Check if getEnv returned empty because the key was missing
		// This check assumes getEnv logs the error if required and missing,
		// but Load should still return a distinct error for a missing required key.
		return Config{}, errors.New("empty key: " + cookieSecureKey)
	}
	c.CookieSecure, err = strconv.ParseBool(cookieSecureStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid value for %s ('%s'): %w", cookieSecureKey, cookieSecureStr, err)
	}

	return c, nil
}

// getEnv retrieves the value of an environment variable.
// If the variable is not set and 'required' is true, it logs an error.
//
// Parameters:
//   - key: The name of the environment variable to retrieve.
//   - required: A boolean indicating if the environment variable is required.
//
// Returns:
//   - string: The value of the environment variable, or an empty string if not set.
func getEnv(key string, required bool) string {
	val := os.Getenv(key)
	if val == "" && required {
		// Use the globally configured logger from the main package or logger package.
		// Avoid reconfiguring the logger here.
		// This log message will use the logger configured in main.
		log.Error().Str("var", key).Msg("Failed to load environment")
	}
	return val
}
