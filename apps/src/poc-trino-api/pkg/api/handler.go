package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/silvioramalho/poc-trino-api/internal/model"
	"github.com/silvioramalho/poc-trino-api/pkg/auth"
	"github.com/silvioramalho/poc-trino-api/pkg/trino"
)

type Handler struct {
	TrinoClient   *trino.Client
	Authenticator *auth.Authenticator
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	catalog := vars["catalog"]
	schema := vars["schema"]
	table := vars["table"]

	// Get query parameters
	queryParams := r.URL.Query()
	startDateStr := queryParams.Get("startDate")
	endDateStr := queryParams.Get("endDate")
	limitStr := queryParams.Get("limit")
	offsetStr := queryParams.Get("offset")

	// Parse parameters
	startDate, _ := time.Parse(time.RFC3339, startDateStr)
	endDate, _ := time.Parse(time.RFC3339, endDateStr)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Create QueryParams
	qParams := model.QueryParams{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
		Offset:    offset,
	}

	query := model.Query{
		Catalog:     catalog,
		Schema:      schema,
		Table:       table,
		QueryParams: qParams,
	}

	data, err := h.TrinoClient.FetchData(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
