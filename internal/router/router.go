package router

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/limiter"
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(rateLimiter *limiter.Limiter) *chi.Mux {
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rateLimiterMiddleware.Limit(next.ServeHTTP).ServeHTTP(w, r)
		})
	})
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})
	return r
}
