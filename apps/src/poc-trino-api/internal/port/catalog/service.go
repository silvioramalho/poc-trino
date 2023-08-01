package catalog

import (
	"github.com/silvioramalho/poc-trino-api/internal/model"
	"github.com/silvioramalho/poc-trino-api/internal/services/trino"
)

type CatalogService interface {
	GetCatalogs() ([]string, error)
	GetSchemas(catalog string) ([]string, error)
	GetTables(catalog, schema string) ([]string, error)
	FetchData(query model.Query) ([]map[string]interface{}, error)
}

type catalog struct {
	Service string
	Uri     string
}

func NewCatalog(service string, uri string) *catalog {
	return &catalog{
		Service: service,
		Uri:     uri,
	}
}

func (c *catalog) GetService() CatalogService {
	switch c.Service {
	case "trino":
		return trino.NewClient(c.Uri)
	default:
		return nil
	}
}
