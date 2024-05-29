package router

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/limiter"
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(rateLimiter *limiter.Limiter) *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.NewRateLimiterMiddleware(rateLimiter).Limit)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Rate Limiter Challenge"))
	})
	return r
}
