package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// New returns a Fiber handler that proxies requests to the target URL.
func New(target string) fiber.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic("invalid proxy target: " + target)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	return adaptor.HTTPHandler(proxy)
}
