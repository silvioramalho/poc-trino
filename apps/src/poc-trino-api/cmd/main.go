package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/silvioramalho/poc-trino-api/internal/config"
	"github.com/silvioramalho/poc-trino-api/internal/handler/auth"
	server "github.com/silvioramalho/poc-trino-api/internal/port/http"
	"github.com/silvioramalho/poc-trino-api/internal/port/trino"
)

func loadEnv() {
	if os.Getenv("ENVIRONMENT") != "k8s" {
		envFile := ".env"
		if os.Getenv("DEBUG") == "true" {
			envFile = "../.env"
		}
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {

	loadEnv()

	cfg := config.LoadConfig()

	authenticator := auth.NewAuthenticator(
		cfg.AuthIssuer,
		cfg.AuthClientID,
		cfg.AuthClientSecret)

	trinoClient := trino.NewClient(cfg.TrinoServerUri)

	server := server.NewServer(authenticator, trinoClient)

	server.Run(cfg.AppServerAddrress)
}
