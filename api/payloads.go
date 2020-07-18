package api

import "github.com/yuriharrison/bexs-test/flightsdb"

// CheapestRoutePayload _
type CheapestRoutePayload struct {
	TotalCost   float64             `json:"total_cost"`
	Description string              `json:"description"`
	Flights     []*flightsdb.Flight `json:"flights"`
}

// FlightPayload _
type FlightPayload struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Price       float64 `json:"price"`
}
