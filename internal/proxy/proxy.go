package proxy

import (
	"net"
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

	// The original director is sufficient if X-User-ID is already set
	// by the RequireAuth middleware on c.Request().Header, which adaptor.HTTPHandler
	// should propagate to the http.Request.
	// proxy.Director remains the default one from NewSingleHostReverseProxy.

	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext, // Added DialTimeout
		ResponseHeaderTimeout: 5 * time.Second,
	}

	return adaptor.HTTPHandler(proxy)
}
