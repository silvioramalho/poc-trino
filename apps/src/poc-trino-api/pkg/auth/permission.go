package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Permissions struct {
	Catalogs map[string]map[string][]string `json:"permissions"`
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

func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		permissions, ok := ctx.Value(permissionsKey).(Permissions)
		if !ok {
			http.Error(w, "Unauthorized: no permissions found", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)

		catalog := vars["catalog"]
		schema := vars["schema"]
		table := vars["table"]

		// O código abaixo é um exemplo e deve ser substituído pela lógica de verificação de permissões real
		if !permissions.CanAccess(catalog, schema, table) {
			http.Error(w, "Unauthorized: insufficient permissions", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
