package initializers

import (
	"log"

	"github.com/lpernett/godotenv"
)

var BasePath string = "upload"

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
