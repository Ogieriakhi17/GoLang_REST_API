package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"

)

type Config struct{
	DatabaseURL string
	Port string
}

func Load() (*Config, error){
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Could not find the environment file")
	}
	var config *Config= &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("PORT"),
	}

	return config, nil
}