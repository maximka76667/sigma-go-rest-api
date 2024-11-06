package main

import (
	"log"

	"github.com/maximka76667/sigma-go-rest-api/database"
)

const (
	PORT = ":8080"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	r := SetupRouter()
	if err := r.RunTLS(PORT, "secrets/cert.pem", "secrets/key.pem"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
