package structs

type Route struct {
	ID          int    `json:"id"`
	Route       string `json:"route"`
	Destination string `json:"destination"`
}
