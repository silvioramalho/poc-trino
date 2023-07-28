package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/silvioramalho/poc-trino-api/config"
	"github.com/silvioramalho/poc-trino-api/pkg/api"
	"github.com/silvioramalho/poc-trino-api/pkg/auth"
	"github.com/silvioramalho/poc-trino-api/pkg/trino"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()

	authenticator := auth.NewAuthenticator(
		cfg.AuthIssuer,
		cfg.AuthClientID,
		cfg.AuthClientSecret)

	trinoClient := trino.NewClient(cfg.TrinoServerUri)

	server := api.NewServer(authenticator, trinoClient)

	server.Run(cfg.AppServerAddrress)
}
