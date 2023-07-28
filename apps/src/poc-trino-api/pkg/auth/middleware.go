package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
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
		rawAccessToken := r.Header.Get("Authorization")
		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}
		rawAccessToken = parts[1]

		idToken, err := a.Verifier.Verify(context.Background(), rawAccessToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to verify ID Token: %v", err), http.StatusUnauthorized)
			return
		}

		var claims jwt.MapClaims
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, fmt.Sprintf("unable to extract claims from ID Token: %v", err), http.StatusUnauthorized)
			return
		}

		permissionsJSON, ok := claims["data_permissions"]
		if !ok {
			http.Error(w, "permissions claim not found in token", http.StatusUnauthorized)
			return
		}

		permissionsArray, ok := permissionsJSON.([]interface{})
		if !ok || len(permissionsArray) < 1 {
			http.Error(w, "permissions claim is not a JSON array or is empty", http.StatusUnauthorized)
			return
		}

		permissionsString, ok := permissionsArray[0].(string)
		if !ok {
			http.Error(w, "permissions claim array does not contain a string", http.StatusUnauthorized)
			return
		}

		var permissions Permissions
		if err := json.Unmarshal([]byte(permissionsString), &permissions); err != nil {
			http.Error(w, "unable to parse permissions claim: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), permissionsKey, permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
