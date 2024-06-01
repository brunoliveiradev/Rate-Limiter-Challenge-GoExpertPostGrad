package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort        int      `mapstructure:"SERVER_PORT"`
	RedisAddr         string   `mapstructure:"REDIS_SERVER_ADDRESS"`
	RedisPassword     string   `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int      `mapstructure:"REDIS_DB"`
	RateLimitByIP     int      `mapstructure:"RATE_LIMIT_BY_IP"`
	RateLimiteByToken int      `mapstructure:"RATE_LIMIT_BY_TOKEN"`
	Tokens            []string `mapstructure:"ALLOWED_TOKENS"`
	Ttl               int      `mapstructure:"TTL_SECONDS"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		ServerPort:        viper.GetInt("SERVER_PORT"),
		RedisAddr:         viper.GetString("REDIS_SERVER_ADDRESS"),
		RedisPassword:     viper.GetString("REDIS_PASSWORD"),
		RedisDB:           viper.GetInt("REDIS_DB"),
		RateLimitByIP:     viper.GetInt("RATE_LIMIT_BY_IP"),
		RateLimiteByToken: viper.GetInt("RATE_LIMIT_BY_TOKEN"),
		Tokens:            viper.GetStringSlice("ALLOWED_TOKENS"),
		Ttl:               viper.GetInt("TTL_SECONDS"),
	}

	return config, nil
}
