package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/rs/zerolog/log"
)

// New returns a Fiber handler that proxies requests to the target URL.
func New(target string) fiber.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Error().Msg("Failed to parse target URL: " + err.Error())
		return func(c *fiber.Ctx) error {
			return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Inject X-User-ID
		if fiberCtx, ok := req.Context().Value("fiber.ctx").(*fiber.Ctx); ok {
			if userID, exists := fiberCtx.Locals("user_id").(string); exists {
				req.Header.Set("X-User-ID", userID)
			}
		}
	}

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	return adaptor.HTTPHandler(proxy)
}
