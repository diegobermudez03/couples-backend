package config

import "os"

type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: getEnv("PORT", ":8081"),
	}
}

func getEnv(envir string, fallCase string) string {
	if val := os.Getenv(envir); val != ""{
		return val 
	}
	return fallCase
}