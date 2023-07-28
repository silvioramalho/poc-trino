package trino

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/silvioramalho/poc-trino-api/internal/model"
	"github.com/trinodb/trino-go-client/trino"
)

type Client struct {
	Host string
}

func NewClient(host string) *Client {
	return &Client{Host: host}
}

// func (c *Client) FetchData(query model.Query) (*model.QueryResult, error) {
func (c *Client) FetchData(query model.Query) (interface{}, error) {

	config := &trino.Config{
		ServerURI:         "http://foobar@my-trino.trino.svc.cluster.local:8080",
		SessionProperties: map[string]string{"query_priority": "1"},
	}

	dsn, err := config.FormatDSN()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("trino", dsn)
	if err != nil {
		log.Println("Error opening connection: ", err)
		return nil, err
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT * FROM %s.%s.%s", query.Catalog, query.Schema, query.Table)
	if query.QueryParams.Offset > 0 {
		sql = fmt.Sprintf("%s ORDER BY name", sql)
	}
	if query.QueryParams.Offset > 0 {
		sql = fmt.Sprintf("%s OFFSET %d", sql, query.QueryParams.Offset)
	}

	if query.QueryParams.Limit > 0 {
		sql = fmt.Sprintf("%s LIMIT %d", sql, query.QueryParams.Limit)
	}

	rows, err := db.Query(sql)
	if err != nil {
		log.Println("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columnsData := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnsData {
			columnPointers[i] = &columnsData[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, columnData := range columnsData {
			rowData[columns[i]] = columnData
		}

		results = append(results, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
