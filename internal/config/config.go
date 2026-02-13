package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"

)

type Config struct{
	DatabaseURL string
	Port string
	JWTSecret string
}

func Load() (*Config, error){
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Could not find the environment file")
	}
	var config *Config= &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return config, nil
}