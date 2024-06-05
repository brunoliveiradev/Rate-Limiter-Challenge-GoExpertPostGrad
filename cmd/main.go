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

	redisClient := persistence.NewRedisClient(config.RedisServerAddr, config.RedisServerPwd, config.RedisDB)
	rateLimiter := limiter.NewLimiter(redisClient, config)

	r := router.NewRouter(rateLimiter)

	log.Printf("Server starting on port %s", config.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), r)
}
