package main

import (
	"github.com/silvioramalho/poc-trino-api/pkg/api"
	"github.com/silvioramalho/poc-trino-api/pkg/auth"
	"github.com/silvioramalho/poc-trino-api/pkg/trino"
)

func main() {
	authenticator := auth.NewAuthenticator(
		"http://keycloak.tools.svc.cluster.local/realms/firehose",
		"trino-api-proxy",
		"e1grRUgS6FZo6tdWtodltwXuhW8Brr6J")

	trinoClient := trino.NewClient("jdbc:trino://my-trino.trino.svc.cluster.local:8080")

	server := api.NewServer(authenticator, trinoClient)

	server.Run(":8080")
}
