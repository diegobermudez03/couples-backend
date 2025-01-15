package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port 		string
	AuthConfig 	*AuthConfig
}

type AuthConfig struct{
	AccessTokenLife 	int64
	JwtSecret 			string
}

func NewConfig() *Config {
	return &Config{
		Port: getEnv("PORT", ":8081"),
		AuthConfig: NewAuthConfig(),
	}
}

func NewAuthConfig() *AuthConfig{
	return &AuthConfig{
		AccessTokenLife: getEnvAsInt64("ACCESS_TOKEN_LIFE", 3600),
		JwtSecret: getEnv("JWT_SECRET", "secret"),
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