package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/silvioramalho/poc-trino-api/pkg/auth"
	"github.com/silvioramalho/poc-trino-api/pkg/trino"
)

type Server struct {
	Router *mux.Router
}

func NewServer(authenticator *auth.Authenticator, trinoClient *trino.Client) *Server {
	handler := &Handler{
		TrinoClient:   trinoClient,
		Authenticator: authenticator,
	}

	r := mux.NewRouter()

	// Apply the authenticator middleware to all routes
	r.Use(authenticator.Middleware)
	r.Use(auth.PermissionMiddleware)

	// Routes
	r.HandleFunc("/{catalog}/{schema}/{table}", handler.Query).Methods("GET")

	return &Server{Router: r}
}

func (s *Server) Run(addr string) {
	result := fmt.Sprintf("Server is running on %v", addr)
	log.Println(result)
	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
