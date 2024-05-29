package configs

import (
	"encoding/json"
	"github.com/spf13/viper"
)

type Config struct {
	RedisAddr             string         `mapstructure:"REDIS_SERVER_ADDRESS"`
	RedisPassword         string         `mapstructure:"REDIS_PASSWORD"`
	RedisDB               int            `mapstructure:"REDIS_DB"`
	RateLimitIP           int            `mapstructure:"RATE_LIMIT_IP"`
	RateLimitTokenDefault int            `mapstructure:"RATE_LIMIT_TOKEN_DEFAULT"`
	TokenAllowed          map[string]int `mapstructure:"RATE_LIMIT_TOKEN_SPECIFIC_ALLOWED"`
	BlockTimeSeconds      int            `mapstructure:"BLOCK_TIME_SECONDS"`
	ServerPort            int            `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		RedisAddr:             viper.GetString("REDIS_SERVER_ADDRESS"),
		RedisPassword:         viper.GetString("REDIS_PASSWORD"),
		RedisDB:               viper.GetInt("REDIS_DB"),
		RateLimitIP:           viper.GetInt("RATE_LIMIT_IP"),
		RateLimitTokenDefault: viper.GetInt("RATE_LIMIT_TOKEN_DEFAULT"),
		BlockTimeSeconds:      viper.GetInt("BLOCK_TIME_SECONDS"),
		ServerPort:            viper.GetInt("SERVER_PORT"),
	}

	specificTokens := viper.GetString("RATE_LIMIT_TOKEN_SPECIFIC_ALLOWED")
	if err := json.Unmarshal([]byte(specificTokens), &config.TokenAllowed); err != nil {
		return nil, err
	}

	return config, nil
}
