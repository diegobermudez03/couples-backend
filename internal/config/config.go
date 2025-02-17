package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port 		string
	AuthConfig 	*AuthConfig
	PostgresConfig *PostgresConfig
	InteractionConfig *InteractionConfig
}

type AuthConfig struct{
	AccessTokenLife 	int64
	RefreshTokenLife 	int64
	JwtSecret 			string
}

type PostgresConfig struct{
	Address 		string
}

type InteractionConfig struct{
	MaxFetchResult	int
}

func NewConfig() *Config {
	return &Config{
		Port: getEnv("PORT", ":8081"),
		AuthConfig: NewAuthConfig(),
		PostgresConfig: NewPostgresConfig(),
		InteractionConfig: NewInteractionConfig(),
	}
}

func NewAuthConfig() *AuthConfig{
	return &AuthConfig{
		AccessTokenLife: getEnvAsInt64("ACCESS_TOKEN_LIFE", 3600),
		JwtSecret: getEnv("JWT_SECRET", "secret"),
		RefreshTokenLife: getEnvAsInt64("REFRESH_TOKEN_LIFE", 1000000000),
	}
}

func NewPostgresConfig() *PostgresConfig{
	return &PostgresConfig{
		Address: getEnv("POSTGRES_DB", ""),
	}
}

func NewInteractionConfig() *InteractionConfig{
	return &InteractionConfig{
		MaxFetchResult: int(getEnvAsInt64("MAX_RESULT_LIMIT", 20)),
	}
}

/////////////////////////////////////////////////

func getEnv(envir string, fallCase string) string {
	if val := os.Getenv(envir); val != ""{
		return val 
	}
	return fallCase
}

func getEnvAsInt64(envir string, fallCase int64) int64{
	unparsed := getEnv(envir, "")
	if unparsed == ""{
		return fallCase
	}
	if num, err := strconv.ParseInt(unparsed, 10, 64); err == nil{
		return num
	}
	return fallCase
}