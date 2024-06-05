package configs

import (
	"encoding/json"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

type Envs struct {
	ServerPort       string   `mapstructure:"SERVER_PORT"`
	RedisServerAddr  string   `mapstructure:"REDIS_SERVER_ADDRESS"`
	RedisServerPwd   string   `mapstructure:"REDIS_SERVER_PASSWORD"`
	RedisDB          int      `mapstructure:"REDIS_DB"`
	RateLimitByIP    int      `mapstructure:"RATE_LIMIT_BY_IP"`
	RateLimitByToken int      `mapstructure:"RATE_LIMIT_BY_TOKEN"`
	AllowedTokens    []string `mapstructure:"ALLOWED_TOKENS"`
	TTLSeconds       int      `mapstructure:"TTL_SECONDS"`
}

func LoadConfig() (*Envs, error) {
	// Tente carregar o arquivo .env
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Sucesso ao ler o arquivo .env
		var cfg Envs
		if err := viper.Unmarshal(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}

	// Falha ao ler o arquivo .env, carregue as variáveis de ambiente
	return loadLocalConfig(), nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func loadLocalConfig() *Envs {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	rateLimitByIP, _ := strconv.Atoi(getEnv("RATE_LIMIT_BY_IP", "2"))
	rateLimitByToken, _ := strconv.Atoi(getEnv("RATE_LIMIT_BY_TOKEN", "10"))
	ttlSeconds, _ := strconv.Atoi(getEnv("TTL_SECONDS", "5"))

	var allowedTokens []string
	if err := json.Unmarshal([]byte(viper.GetString("ALLOWED_TOKENS")), &allowedTokens); err != nil {
		allowedTokens = []string{"token1", "token2", "token3"}
	}

	return &Envs{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		RedisServerAddr:  getEnv("REDIS_SERVER_ADDRESS", "localhost:6379"),
		RedisServerPwd:   getEnv("REDIS_SERVER_PASSWORD", ""),
		RedisDB:          redisDB,
		RateLimitByIP:    rateLimitByIP,
		RateLimitByToken: rateLimitByToken,
		AllowedTokens:    allowedTokens,
		TTLSeconds:       ttlSeconds,
	}
}
