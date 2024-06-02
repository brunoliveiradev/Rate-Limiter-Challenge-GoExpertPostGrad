package main

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/configs"
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/limiter"
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/router"
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/persistence"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	redisClient := persistence.NewRedisClient(config.RedisAddr, config.RedisPassword, config.RedisDB)
	rateLimiter := limiter.NewLimiter(redisClient, config)

	r := router.NewRouter(rateLimiter)

	log.Printf("Server starting on port %d", config.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), r)
}
