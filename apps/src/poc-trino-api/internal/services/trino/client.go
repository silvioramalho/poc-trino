package trino

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/silvioramalho/poc-trino-api/internal/model"
	"github.com/trinodb/trino-go-client/trino"
)

type Client struct {
	ServerURI string
}

func NewClient(serverURI string) *Client {
	return &Client{ServerURI: serverURI}
}

func (c *Client) createTrinoConn() (*sql.DB, error) {
	config := &trino.Config{
		ServerURI:         c.ServerURI,
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

	return db, nil
}

func (c *Client) buildSQL(query model.Query) string {
	sql := fmt.Sprintf("SELECT * FROM %s.%s.%s", query.Catalog, query.Schema, query.Table)
	// ATTENTION: This is just a POC. Forcing using name field for using OFFSET.
	if query.QueryParams.Offset > 0 {
		sql = fmt.Sprintf("%s ORDER BY name", sql)
	}
	if query.QueryParams.Offset > 0 {
		sql = fmt.Sprintf("%s OFFSET %d", sql, query.QueryParams.Offset)
	}

	if query.QueryParams.Limit > 0 {
		sql = fmt.Sprintf("%s LIMIT %d", sql, query.QueryParams.Limit)
	}

	return sql
}

func processRows(rows *sql.Rows) ([]map[string]interface{}, error) {
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

func (c *Client) FetchData(query model.Query) ([]map[string]interface{}, error) {
	db, err := c.createTrinoConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := c.buildSQL(query)

	rows, err := db.Query(sql)
	if err != nil {
		log.Println("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	return processRows(rows)
}

func (c *Client) GetCatalogs() ([]string, error) {
	db, err := c.createTrinoConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SHOW CATALOGS")
	if err != nil {
		log.Println("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	var catalogs []string
	for rows.Next() {
		var catalog string
		if err := rows.Scan(&catalog); err != nil {
			return nil, err
		}
		catalogs = append(catalogs, catalog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return catalogs, nil
}

func (c *Client) GetSchemas(catalog string) ([]string, error) {
	db, err := c.createTrinoConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := fmt.Sprintf("SHOW SCHEMAS FROM %s", catalog)
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

func (c *Client) GetTables(catalog, schema string) ([]string, error) {
	db, err := c.createTrinoConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := fmt.Sprintf("SHOW TABLES FROM %s.%s", catalog, schema)
	rows, err := db.Query(sql)
	if err != nil {
		log.Println("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
