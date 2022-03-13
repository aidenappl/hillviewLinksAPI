package structs

import "time"

type Route struct {
	ID          int       `json:"id"`
	Route       string    `json:"route"`
	Destination string    `json:"destination"`
	CreatedBy   int       `json:"created_by"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}
