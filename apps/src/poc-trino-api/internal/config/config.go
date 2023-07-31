package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AuthIssuer        string `envconfig:"AUTH_ISSUER" required:"true"`
	AuthClientID      string `envconfig:"AUTH_CLIENT_ID" required:"true"`
	AuthClientSecret  string `envconfig:"AUTH_CLIENT_SECRET" required:"true"`
	TrinoServerUri    string `envconfig:"TRINO_SERVER_URI" required:"true"`
	AppServerAddrress string `envconfig:"SERVER_ADDR" required:"true"`
}

func LoadConfig() *Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
