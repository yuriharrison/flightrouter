package loader

import (
	"testing"

	flightsdb "github.com/yuriharrison/bexs-test/flightsdb"
)

const (
	flightsFile = "../fixtures/flights.csv"
	numFlights  = 10 // 10 flights 3 updates
)

func TestImportFlightsFromFile(t *testing.T) {
	db := flightsdb.New()
	err := ImportFlightsFromFile(flightsFile, db)
	if err != nil {
		t.Error(err.Error())
	}
	if db.Size() != numFlights {
		t.Errorf("Size %v != %v numFlights", db.Size(), numFlights)
	}
}
