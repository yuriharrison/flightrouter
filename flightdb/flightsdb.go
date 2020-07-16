package flightsdb

import (
	"fmt"
	"strings"
)

// Airport struct
type Airport struct {
	Code    string
	flights map[*Airport]*Flight
}

// NewAirport create a new *FlightsDB
func NewAirport(code string) *Airport {
	return &Airport{
		code,
		make(map[*Airport]*Flight),
	}
}

// Flight struct
type Flight struct {
	orig  *Airport
	dest  *Airport
	Price float32
}

// NewFlight create flight or update price for destination
// return true if a new flight was created false if updated.
func (ap *Airport) NewFlight(dest *Airport, price float32) bool {
	if flight, ok := ap.flights[dest]; ok {
		flight.Price = price
		return false
	}
	ap.flights[dest] = &Flight{ap, dest, price}
	return true
}

// FlightsDB database for query flights
type FlightsDB struct {
	airports   map[string]*Airport
	numFlights int
	cache      Cache
}

// NewFlightsDB create a new *FlightsDB
func NewFlightsDB() *FlightsDB {
	return &FlightsDB{
		airports: make(map[string]*Airport),
		cache:    Cache{},
	}
}

func (db *FlightsDB) getAirports(origCode, destCode string) (*Airport, *Airport, error) {
	origin, err := db.GetAirport(origCode)
	if err != nil {
		return nil, nil, err
	}
	destination, err := db.GetAirport(destCode)
	if err != nil {
		return nil, nil, err
	}
	return origin, destination, nil
}

// Add new flight to DB
func (db *FlightsDB) Add(origCode, destCode string, price float32) error {
	origin, destination, err := db.getAirports(origCode, destCode)
	if err != nil {
		return err
	}
	db.cache.Clean()
	if new := origin.NewFlight(destination, price); new {
		db.numFlights++
	}
	return nil
}

// GetAirport return the existent Airport or create a new one for [code]
func (db *FlightsDB) GetAirport(code string) (*Airport, error) {
	if len(code) != 3 {
		return nil, fmt.Errorf(
			"Invalid airport code format. Expecting exacly 3 characters got %v",
			code,
		)
	}
	code = strings.ToUpper(code)

	if ap, ok := db.airports[code]; ok {
		return ap, nil
	}
	newAirport := NewAirport(code)
	db.airports[code] = newAirport
	return newAirport, nil
}

// Size return the total number of flights on the database
func (db *FlightsDB) Size() int {
	return db.numFlights
}

// Remove removes a flight if it exists
func (db *FlightsDB) Remove(origCode, destCode string) error {
	origin, destination, err := db.getAirports(origCode, destCode)
	if err != nil {
		return err
	}
	delete(origin.flights, destination)
	db.cache.Clean()
	return nil
}

// CheapestRoute return the cheapest route for a given origem and destination
func (db *FlightsDB) CheapestRoute(origCode, destCode string) ([]*Flight, error) {
	if cache := db.cache.GetCheapestRoute(origCode, destCode); cache != nil {
		return cache, nil
	}
	origin, destination, err := db.getAirports(origCode, destCode)
	if err != nil {
		return nil, err
	}
	route := FindCheapestRoute(origin, destination)
	if route == nil {
		return nil, fmt.Errorf(
			"There is no route available for %v-%v",
			origCode,
			destCode,
		)
	}
	db.cache.SetCheapestRoute(origCode, destCode, route)
	return route, nil
}
