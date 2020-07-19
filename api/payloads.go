package api

import (
	"time"

	"github.com/yuriharrison/flightrouter/flightsdb"
)

// CheapestRoutePayload _
type CheapestRoutePayload struct {
	TotalCost   float64             `json:"totalCost"`
	Description string              `json:"description"`
	Flights     []*flightsdb.Flight `json:"flights"`
}

// FlightPayload _
type FlightPayload struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Price       float64 `json:"price"`
}

// StatusCache _
type StatusCache struct {
	Hits   int `json:"hits"`
	Misses int `json:"misses"`
}

// StatusPayload _
type StatusPayload struct {
	NumFlights int         `json:"numberFlights"`
	Cache      StatusCache `json:"cache"`
	Now        time.Time   `json:"time"`
}
