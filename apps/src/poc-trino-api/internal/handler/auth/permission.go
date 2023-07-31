package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type Permissions struct {
	Catalogs map[string]map[string][]string `json:"permissions"`
}

func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permissions, err := getPermissionsFromContext(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		catalog, schema, table := getCatalogSchemaTable(mux.Vars(r))

		// O código abaixo é um exemplo e deve ser substituído pela lógica de verificação de permissões real
		if !permissions.CanAccess(catalog, schema, table) {
			http.Error(w, "Unauthorized: insufficient permissions", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p *Permissions) CanAccess(catalog, schema, table string) bool {
	// Verifique se o catálogo existe.
	if schemas, ok := p.Catalogs[catalog]; ok {
		// O catálogo existe, verifique se o esquema existe.
		if tables, ok := schemas[schema]; ok {
			// O esquema existe, verifique se a tabela existe.
			for _, t := range tables {
				if t == table {
					return true
				}
			}
		}
	}

	return false
}

func getPermissionsFromContext(r *http.Request) (Permissions, error) {
	ctx := r.Context()
	permissions, ok := ctx.Value(permissionsKey).(Permissions)
	if !ok {
		return Permissions{}, errors.New("no permissions found")
	}

	return permissions, nil
}

func getCatalogSchemaTable(vars map[string]string) (string, string, string) {
	catalog := vars["catalog"]
	schema := vars["schema"]
	table := vars["table"]

	return catalog, schema, table
}

func extractPermissions(idToken *oidc.IDToken) (Permissions, error) {
	var claims jwt.MapClaims
	if err := idToken.Claims(&claims); err != nil {
		return Permissions{}, fmt.Errorf("unable to extract claims from ID Token: %w", err)
	}

	permissionsJSON, ok := claims["data_permissions"]
	if !ok {
		return Permissions{}, errors.New("permissions claim not found in token")
	}

	permissionsArray, ok := permissionsJSON.([]interface{})
	if !ok || len(permissionsArray) < 1 {
		return Permissions{}, errors.New("permissions claim is not a JSON array or is empty")
	}

	permissionsString, ok := permissionsArray[0].(string)
	if !ok {
		return Permissions{}, errors.New("permissions claim array does not contain a string")
	}

	var permissions Permissions
	if err := json.Unmarshal([]byte(permissionsString), &permissions); err != nil {
		return Permissions{}, fmt.Errorf("unable to parse permissions claim: %w", err)
	}

	return permissions, nil
}
