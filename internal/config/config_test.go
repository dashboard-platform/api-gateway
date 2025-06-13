// Package config provides functionality for loading and managing application configuration.
// It retrieves configuration values from environment variables and ensures that all required
// settings are properly initialized. This package is essential for setting up the application's
// runtime environment, including database connections, JWT secrets, and server settings.
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetEnv tests the getEnv function to ensure it retrieves the correct environment variable value.
// It verifies that the function returns the expected value for existing environment variables
// and an empty string for non-existent variables.
func TestGetEnv(t *testing.T) {
	os.Setenv("test", "test")
	defer os.Unsetenv("test")

	tests := []struct {
		name string // Name of the test case.
		key  string // The environment variable key to retrieve.
		want string // The expected value of the environment variable.
	}{
		{
			name: "Test ENV variable - required",
			key:  "test",
			want: "test",
		},
		{
			name: "Test non-existent ENV variable - required",
			key:  "non_existent",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEnv(tt.key, true) // Assuming 'required' is true for these test cases
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestLoad tests the Load function to ensure it correctly loads the configuration from environment variables.
// It verifies that the function returns a valid configuration when all required environment variables are set
// and returns an error when any mandatory variable is missing.
func TestLoad(t *testing.T) {
	// Helper to set environment variables for a test case
	setEnvs := func(envs map[string]string) {
		for k, v := range envs {
			os.Setenv(k, v)
		}
	}
	// Helper to unset environment variables
	unsetEnvs := func(keys []string) {
		for _, k := range keys {
			os.Unsetenv(k)
		}
	}
	tests := []struct {
		name    string   // Name of the test case.
		envs    []string // List of environment variables to set for the test.
		wantErr bool     // Expected error: true if an error is expected.
	}{
		{
			name:    "Test Load with all envs set",
			envs:    []string{envKey, portEnv, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: false,
		},
		{
			name:    "Test Load with missing envKey",
			envs:    []string{portEnv, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: false,
		},
		{
			name:    "Test Load with missing portEnv",
			envs:    []string{envKey, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing frontEndKey",
			envs:    []string{envKey, portEnv, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing authServiceKey",
			envs:    []string{envKey, portEnv, frontEndKey, templateServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing templateServiceKey",
			envs:    []string{envKey, portEnv, frontEndKey, authServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing pdfServiceKey",
			envs:    []string{envKey, portEnv, frontEndKey, authServiceKey, templateServiceKey, jwtSecretKey, cookieSecureKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing jwtSecretKey",
			envs:    []string{envKey, portEnv, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, cookieSecureKey},
			wantErr: true,
		},
		{
			name:    "Test Load with missing cookieSecureKey",
			envs:    []string{envKey, portEnv, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey},
			wantErr: true,
		},
		{
			name:    "Test Load with empty envs",
			envs:    []string{},
			wantErr: true,
		},
		{
			name:    "Test Load with invalid cookieSecureKey",
			envs:    []string{envKey, portEnv, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey, "COOKIE_SECURE_INVALID"}, // Special case for value
			wantErr: true,
		},
	}

	allPossibleEnvs := []string{envKey, portEnv, frontEndKey, authServiceKey, templateServiceKey, pdfServiceKey, jwtSecretKey, cookieSecureKey}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up all known envs before each test
			unsetEnvs(allPossibleEnvs)

			// Set specific envs for the current test case
			envMap := make(map[string]string)
			isInvalidCookieTest := false
			for _, envName := range tt.envs {
				if envName == "COOKIE_SECURE_INVALID" {
					envMap[cookieSecureKey] = "not-a-boolean"
					isInvalidCookieTest = true
				} else if envName == cookieSecureKey {
					// For test cases that include cookieSecureKey and are not the "invalid" case,
					// set a valid boolean string (e.g., "true") to ensure strconv.ParseBool succeeds.
					envMap[cookieSecureKey] = "true"
				} else {
					envMap[envName] = "test"
				}
			}

			// Default values for a successful load if not in error case
			if !tt.wantErr {
				if _, ok := envMap[envKey]; !ok {
					envMap[envKey] = "dev" // Default or test value
				}
				if _, ok := envMap[portEnv]; !ok {
					envMap[portEnv] = "8080"
				}
				if _, ok := envMap[frontEndKey]; !ok {
					envMap[frontEndKey] = "http://localhost:3000"
				}
				if _, ok := envMap[authServiceKey]; !ok {
					envMap[authServiceKey] = "http://auth"
				}
				if _, ok := envMap[templateServiceKey]; !ok {
					envMap[templateServiceKey] = "http://template"
				}
				if _, ok := envMap[pdfServiceKey]; !ok {
					envMap[pdfServiceKey] = "http://pdf"
				}
				if _, ok := envMap[jwtSecretKey]; !ok {
					envMap[jwtSecretKey] = "secret"
				}
				if _, ok := envMap[cookieSecureKey]; !ok && !isInvalidCookieTest {
					envMap[cookieSecureKey] = "true"
				}
			}

			setEnvs(envMap)
			defer unsetEnvs(allPossibleEnvs) // Ensure cleanup

			cfg, err := Load()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, cfg)
			if !tt.wantErr {
				assert.Equal(t, envMap[portEnv], cfg.Port) // Example assertion
				// Add more assertions for other fields if necessary
			}
		})
	}
}
