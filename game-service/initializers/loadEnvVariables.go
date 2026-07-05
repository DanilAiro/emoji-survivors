package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	_ = godotenv.Load()
}
