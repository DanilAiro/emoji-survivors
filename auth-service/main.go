package authservice

import (
	"emoji-survivors/auth-service/initializers"
	"emoji-survivors/auth-service/repository"
)

func main() {
	initializers.LoadEnvVariables()

	repository.ConnectDB()

	
}