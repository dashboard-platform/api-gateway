package proxy

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type responseRecorder struct {
	ctx *fiber.Ctx
}

func (r *responseRecorder) Header() http.Header {
	h := make(http.Header)
	r.ctx.Response().Header.VisitAll(func(k, v []byte) {
		h.Set(string(k), string(v))
	})
	return h
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.ctx.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.ctx.Status(statusCode)
}
