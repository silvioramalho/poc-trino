package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Verifier *oidc.IDTokenVerifier
	Config   *oauth2.Config
}

type key int

const (
	permissionsKey key = iota
)

func NewAuthenticator(issuer, clientID, clientSecret string) *Authenticator {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		panic(err)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID:        clientID,
		SkipIssuerCheck: true,
	}

	return &Authenticator{
		Verifier: provider.Verifier(oidcConfig),
		Config:   config,
	}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract raw access token from the Authorization header
		rawAccessToken, err := extractToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify the access token
		idToken, err := a.Verifier.Verify(context.Background(), rawAccessToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to verify ID Token: %v", err), http.StatusUnauthorized)
			return
		}

		// Extract permissions from the ID Token
		permissions, err := extractPermissions(idToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to extract permissions: %v", err), http.StatusUnauthorized)
			return
		}

		// Add permissions to context and call the next handler
		ctx := context.WithValue(r.Context(), permissionsKey, permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	rawAccessToken := r.Header.Get("Authorization")
	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid Authorization header")
	}
	return parts[1], nil
}
