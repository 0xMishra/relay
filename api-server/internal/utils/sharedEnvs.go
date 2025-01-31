package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// loading the .env file
var err = godotenv.Load()

// random PROJECT ID
var (
	RedisUrl = os.Getenv("REDIS_URL")
)
