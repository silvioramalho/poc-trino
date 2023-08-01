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

type permissions struct {
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

		if !permissions.CanAccess(catalog, schema, table) {
			http.Error(w, "Unauthorized: insufficient permissions", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p *permissions) CanAccess(catalog, schema, table string) bool {
	if schemas, ok := p.Catalogs[catalog]; ok {
		if tables, ok := schemas[schema]; ok {
			for _, t := range tables {
				if t == table {
					return true
				}
			}
		}
	}

	return false
}

func getPermissionsFromContext(r *http.Request) (permissions, error) {
	ctx := r.Context()
	p, ok := ctx.Value(permissionsKey).(permissions)
	if !ok {
		return permissions{}, errors.New("no permissions found")
	}

	return p, nil
}

func getCatalogSchemaTable(vars map[string]string) (string, string, string) {
	catalog := vars["catalog"]
	schema := vars["schema"]
	table := vars["table"]

	return catalog, schema, table
}

func extractPermissions(idToken *oidc.IDToken) (permissions, error) {
	var claims jwt.MapClaims
	if err := idToken.Claims(&claims); err != nil {
		return permissions{}, fmt.Errorf("unable to extract claims from ID Token: %w", err)
	}

	permissionsJSON, ok := claims["data_permissions"]
	if !ok {
		return permissions{}, errors.New("permissions claim not found in token")
	}

	permissionsArray, ok := permissionsJSON.([]interface{})
	if !ok || len(permissionsArray) < 1 {
		return permissions{}, errors.New("permissions claim is not a JSON array or is empty")
	}

	permissionsString, ok := permissionsArray[0].(string)
	if !ok {
		return permissions{}, errors.New("permissions claim array does not contain a string")
	}

	var p permissions
	if err := json.Unmarshal([]byte(permissionsString), &p); err != nil {
		return permissions{}, fmt.Errorf("unable to parse permissions claim: %w", err)
	}

	return p, nil
}
