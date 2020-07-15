package flightsdb

import (
	"container/heap"
	"fmt"
	"strings"
)

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

type Flight struct {
	orig  *Airport
	dest  *Airport
	Price float32
}

// NewFlight create flight or update price for destination
// return true if a new flight was created false if updated.
// *It was not espeficied what is the expected behavior when
// the same flight is inserted more than once, thus I'll be
// updating the flight price value when that happens.
func (ap *Airport) NewFlight(dest *Airport, price float32) bool {
	if flight, ok := ap.flights[dest]; ok {
		flight.Price = price
		return false
	}
	ap.flights[dest] = &Flight{ap, dest, price}
	return true
}

type FlightsDB struct {
	airports   map[string]*Airport
	numFlights int
}

// NewFlightsDB create a new *FlightsDB
func NewFlightsDB() *FlightsDB {
	return &FlightsDB{
		airports: make(map[string]*Airport),
	}
}

// Add new flight to DB
func (fdb *FlightsDB) Add(originCode, destCode string, price float32) {
	origin, destination := fdb.GetAirport(originCode), fdb.GetAirport(destCode)
	if new := origin.NewFlight(destination, price); new {
		fdb.numFlights++
	}
}

// GetAirport return the existent Airport or create a new one for [code]
func (fdb *FlightsDB) GetAirport(code string) *Airport {
	code = strings.ToUpper(code)
	if ap, ok := fdb.airports[code]; ok {
		return ap
	}
	newAirport := NewAirport(code)
	fdb.airports[code] = newAirport
	return newAirport
}

// Size return the total number of flights on the database
func (fdb *FlightsDB) Size() int {
	return fdb.numFlights
}

// Remove removes a flight if it exists
func (fdb *FlightsDB) Remove(origCode, destcode string) {
	origin, destination := fdb.GetAirport(origCode), fdb.GetAirport(destcode)
	delete(origin.flights, destination)
}

// CheapestRoute return the cheapest route for a given origem and destination
func (fdb *FlightsDB) CheapestRoute(origCode, destCode string) ([]*Flight, error) {
	origin, destination := fdb.GetAirport(origCode), fdb.GetAirport(destCode)
	visited := make(map[*Airport]bool)
	accPriceTable := make(map[*Airport]float32)
	routeTrace := make(map[*Airport]*Flight)

	queue := &Queue{&QueueItem{data: origin, value: 0, index: 0}}
	accPriceTable[origin] = 0
	heap.Init(queue)
	for queue.Len() > 0 {
		if _, ok := visited[destination]; ok {
			break
		}

		item := heap.Pop(queue).(*QueueItem)
		ap := item.data
		for _, flight := range ap.flights {
			dest := flight.dest
			accPrice := accPriceTable[ap] + flight.Price
			if oldAccPrice, ok := accPriceTable[dest]; !ok || accPrice < oldAccPrice {
				accPriceTable[dest] = accPrice
				routeTrace[dest] = flight
				heap.Push(queue, &QueueItem{data: dest, value: accPrice})
			}
		}
		visited[ap] = true
	}

	if _, ok := visited[destination]; !ok {
		return nil, fmt.Errorf(
			"There is no route available for %v-%v",
			origCode,
			destCode,
		)
	}

	path := []*Flight{}
	for next := routeTrace[destination]; next != nil; next = routeTrace[next.orig] {
		path = append(path, next)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, nil
}
