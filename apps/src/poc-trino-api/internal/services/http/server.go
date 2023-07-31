package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/silvioramalho/poc-trino-api/internal/handler/auth"
	handler "github.com/silvioramalho/poc-trino-api/internal/handler/http"
	"github.com/silvioramalho/poc-trino-api/internal/services/trino"
)

type Server struct {
	Router *mux.Router
}

func NewServer(authenticator *auth.Authenticator, trinoClient *trino.Client) *Server {
	handler := &handler.Handler{
		TrinoClient:   trinoClient,
		Authenticator: authenticator,
	}

	r := mux.NewRouter()
	r.Use(authenticator.Middleware)
	r.Use(auth.PermissionMiddleware)
	r.Use(handler.ContentTypeMiddleware)

	// Routes
	r.HandleFunc("/{catalog}/{schema}/{table}", handler.Query).Methods("GET")

	return &Server{Router: r}
}

func (s *Server) Run(addr string) {
	log.Printf("Server is starting on %v\n", addr)
	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		log.Fatalf("Failed to start server on %v: %v\n", addr, err)
	}
}
