package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/silvioramalho/poc-trino-api/internal/handler/auth"
	"github.com/silvioramalho/poc-trino-api/internal/model"
	"github.com/silvioramalho/poc-trino-api/internal/port/trino"
)

type Handler struct {
	TrinoClient   *trino.Client
	Authenticator *auth.Authenticator
}

func getRouteVars(vars map[string]string) (string, string, string, error) {
	catalog, ok := vars["catalog"]
	if !ok {
		return "", "", "", errors.New("missing route variable: catalog")
	}

	schema, ok := vars["schema"]
	if !ok {
		return "", "", "", errors.New("missing route variable: schema")
	}

	table, ok := vars["table"]
	if !ok {
		return "", "", "", errors.New("missing route variable: table")
	}

	return catalog, schema, table, nil
}

func parseQueryParams(queryParams url.Values) (model.QueryParams, error) {

	var qParams model.QueryParams

	startDateStr := queryParams.Get("startDate")
	if startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return qParams, fmt.Errorf("invalid startDate: %w", err)
		}
		qParams.StartDate = startDate
	}

	endDateStr := queryParams.Get("endDate")
	if endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return qParams, fmt.Errorf("invalid endDate: %w", err)
		}
		qParams.EndDate = endDate
	}

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return qParams, fmt.Errorf("invalid limit: %w", err)
		}
		qParams.Limit = limit
	}

	offsetStr := queryParams.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return qParams, fmt.Errorf("invalid offset: %w", err)
		}
		qParams.Offset = offset
	}

	return qParams, nil
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	catalog, schema, table, err := getRouteVars(mux.Vars(r))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error extracting route variables: %v", err), http.StatusBadRequest)
		return
	}

	qParams, err := parseQueryParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing query parameters: %v", err), http.StatusBadRequest)
		return
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

	json.NewEncoder(w).Encode(data)
}
