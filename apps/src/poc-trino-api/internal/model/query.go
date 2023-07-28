package model

type Query struct {
	Catalog     string
	Schema      string
	Table       string
	QueryParams QueryParams
}
