package app

import (
	"Auth/internal/controller"
	"github.com/joho/godotenv"
	"log"
)

func Application() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	controller.Registry()
}
