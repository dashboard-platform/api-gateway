package main

import (
	"context"
	"os"
	"os/signal"
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
		log.Fatal().Err(err).Msg("Failed to load configuration")
		// os.Exit(1) // log.Fatal() already exits with status 1
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

		// Add custom request logger middleware.
		middleware.RequestLogger(httpLogger),
	)

	// Proxy handlers
	authProxy := proxy.New(c.AuthServiceURL)
	templatesProxy := proxy.New(c.TemplateServiceURL)
	pdfProxy := proxy.New(c.PDFServiceURL)

	// JWT object for authentication middleware
	jwtObj := &middleware.JWTObj{
		Secret: c.JWTSecret,
	}

	globalLimiter := limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
	})

	// Routes
	app.All("/auth/*",
		globalLimiter,
		authProxy,
	)
	app.Post("/templates/:id/preview",
		middleware.RequireAuth(jwtObj),
		limiter.New(limiter.Config{
			Max:        1000,
			Expiration: 1 * time.Minute,
		}),
		templatesProxy,
	)
	app.All("/templates/*",
		middleware.RequireAuth(jwtObj),
		globalLimiter,
		templatesProxy,
	)
	app.All("/pdf/*",
		middleware.RequireAuth(jwtObj),
		globalLimiter,
		pdfProxy,
	)

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("api-gateway is alive")
	})
	app.Get("/logout", func(ctx *fiber.Ctx) error {
		ctx.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			Secure:   c.CookieSecure,
			HTTPOnly: true,
			SameSite: "None",
		})
		return ctx.SendStatus(fiber.StatusOK)
	})

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // syscall.SIGINT, syscall.SIGTERM

	// Goroutine to start the server
	go func() {
		log.Info().Msgf("API Gateway starting on %s", c.Port)
		if err := app.Listen(c.Port); err != nil {
			log.Error().Err(err).Msg("Error starting API gateway")
			quit <- os.Interrupt // Signal main to exit if server fails to start
		}
	}()

	// Wait for an OS signal
	<-quit
	log.Info().Msg("Shutting down API Gateway...")

	// Attempt to gracefully shut down the server
	if err := app.ShutdownWithContext(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	}
	log.Info().Msg("API Gateway gracefully stopped")
}
