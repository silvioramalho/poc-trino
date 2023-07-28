package model

import "time"

type QueryParams struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}
