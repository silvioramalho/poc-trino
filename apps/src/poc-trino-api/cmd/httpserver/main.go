package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/silvioramalho/poc-trino-api/internal/config"
	"github.com/silvioramalho/poc-trino-api/internal/handler/auth"
	handler "github.com/silvioramalho/poc-trino-api/internal/handler/http"
	"github.com/silvioramalho/poc-trino-api/internal/port/catalog"
	server "github.com/silvioramalho/poc-trino-api/internal/services/http"
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

	catalogCfg := catalog.NewCatalog(cfg.CatalogService, cfg.CatalogServerUri)
	catalogService := catalogCfg.GetService()

	authenticator := auth.NewAuthenticator(
		cfg.AuthIssuer,
		cfg.AuthClientID,
		cfg.AuthClientSecret)

	handler := handler.NewHandler(catalogService)

	server := server.NewServer(handler, authenticator)

	server.Run(cfg.AppServerAddrress)
}
