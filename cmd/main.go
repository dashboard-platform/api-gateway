package main

import (
	"time"

	"github.com/dashboard-platform/api-gateway/internal/config"
	"github.com/dashboard-platform/api-gateway/internal/logger"
	"github.com/dashboard-platform/api-gateway/internal/middleware"
	"github.com/dashboard-platform/api-gateway/internal/proxy"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	// Load the configuration from environment variables.
	c, err := config.Load()
	if err != nil {
		return
	}

	// Initialize the logger with the loaded configuration
	baseLogger := logger.Init(c.Env)
	httpLogger := logger.NewComponentLogger(baseLogger, "http")

	app := fiber.New()
	// Middlewares
	app.Use(
		cors.New(cors.Config{
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			AllowMethods:     "GET, POST, PUT, DELETE",
			AllowOrigins:     c.FrontendURL,
			AllowCredentials: true,
		}),

		helmet.New(),

		//csrf.New(),

		limiter.New(limiter.Config{
			Max:        20,
			Expiration: 1 * time.Minute,
		}),

		// Add custom request logger middleware.
		middleware.RequestLogger(httpLogger),
	)

	// Proxy handlers
	authProxy := proxy.New(c.AuthServiceURL)
	templatesProxy := proxy.New(c.TemplateServiceURL)

	// JWT object for authentication middleware
	jwtObj := &middleware.JWTObj{
		Secret: []byte(c.JWTSecret),
	}

	// Routes
	app.All("/auth/*", authProxy)
	app.All("/templates/*", middleware.RequireAuth(jwtObj), templatesProxy)

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("api-gateway is alive")
	})

	// Start the HTTP server.
	log.Info().Msgf("API Gateway started on %s", c.Port)
	if err = app.Listen(c.Port); err != nil {
		log.Error().Msgf("Error starting api gateway: %v", err)
		return
	}
}
