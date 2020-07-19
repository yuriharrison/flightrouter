package flightsdb

import (
	"fmt"
	"strings"
	"testing"
)

const (
	APQ = "APQ"
	BAZ = "BAZ"
	BEL = "BEL"
	BSB = "BSB"
	CNF = "CNF"
	GRU = "GRU"
	PLU = "PLU"
	VCP = "VCP"
)

var dataSize int = 10
var data = []struct {
	orig, dest string
	price      float64
}{
	{CNF, GRU, 200},
	{CNF, GRU, 215.5},
	{CNF, APQ, 899.99},
	{GRU, CNF, 70},
	{GRU, CNF, 78},
	{GRU, BAZ, 300},
	{GRU, PLU, 10},
	{BAZ, BEL, 600},
	{BAZ, GRU, 10},
	{BAZ, PLU, 200},
	{PLU, APQ, 15},
	{BSB, VCP, 415.5},
}

type assertRoutePayload struct {
	t                  *testing.T
	db                 *FlightsDB
	origCode, destCode string
	expectedRoute      []string
}

func assertRoute(p assertRoutePayload) {
	route, err := p.db.CheapestRoute(p.origCode, p.destCode)
	if err != nil {
		p.t.Error(err)
	}

	fail := false
	if len(p.expectedRoute) != len(route)+1 {
		fail = true
	}
	var routeCodes []string
	for i, flight := range route {
		if i == 0 {
			routeCodes = append(routeCodes, flight.Origin.Code)
			if !fail && flight.Origin.Code != strings.ToUpper(p.expectedRoute[i]) {
				fail = true
			}
		}
		routeCodes = append(routeCodes, flight.Destination.Code)
		if !fail && flight.Destination.Code != strings.ToUpper(p.expectedRoute[i+1]) {
			fail = true
		}
	}
	if fail {
		makeRouteText := func(codes []string) string {
			text := ""
			for i, r := range codes {
				if i < 1 {
					text += r
					continue
				}
				text += "-" + r
			}
			return text
		}
		fmt.Println("Expected:\t", makeRouteText(p.expectedRoute))
		fmt.Println("Got:\t\t", makeRouteText(routeCodes))
		p.t.Errorf("Routes don't match")
	}
}

func TestFlightsDB(t *testing.T) {
	db := New()
	for _, d := range data {
		err := db.Add(d.orig, d.dest, d.price)
		if err != nil {
			t.Errorf("Error adding flight %v %v %v", d.orig, d.dest, d.price)
		}
	}

	// Assert size
	if l := db.Size(); l != dataSize {
		t.Errorf("FlightsDB.Size() %v != %v expected", l, 1)
	}

	// Assert cheapest route
	assertRoute(assertRoutePayload{
		t, db, BAZ, BEL, []string{BAZ, BEL},
	})
	assertRoute(assertRoutePayload{
		t, db, CNF, APQ, []string{CNF, GRU, PLU, APQ},
	})
	db.Remove(GRU, PLU)
	assertRoute(assertRoutePayload{
		t, db, CNF, APQ, []string{CNF, GRU, BAZ, PLU, APQ},
	})
	db.Add(BAZ, PLU, 400)
	assertRoute(assertRoutePayload{
		t, db, CNF, APQ, []string{CNF, APQ},
	})

	// Assert NO routes available
	if route, err := db.CheapestRoute(CNF, BSB); err == nil || len(route) > 0 {
		t.Error(err)
	}
	if route, err := db.CheapestRoute("AAA", "ZZZ"); err == nil || len(route) > 0 {
		t.Error(err)
	}

	// Remove invalid, don't crash
	db.Remove("AAA", "ZZZ")
}

func TestFlightDBFailInvalidCodeFormat(t *testing.T) {
	db := New()
	if err := db.Add("AAA4", "ZZZ", 200); err == nil {
		t.Error("Should fail .Add.")
	}
	if err := db.Add("AAA", "ZZ", 200); err == nil {
		t.Error("Should fail .Add.")
	}
	if err := db.Add("AA", "ZZZ", 200); err == nil {
		t.Error("Should fail .Add.")
	}
	if err := db.Remove("AA", "ZZZ"); err == nil {
		t.Error("Should fail .Remove.")
	}
	if _, err := db.CheapestRoute("AA", "ZZ"); err == nil {
		t.Error("Should fail .CheapestRoute.")
	}
}

func TestFlightsDBCache(t *testing.T) {
	db := New()
	db.Add(CNF, BSB, 200)
	db.Add(CNF, GRU, 200)
	db.CheapestRoute(CNF, BSB)                  // miss
	db.CheapestRoute(CNF, GRU)                  // miss
	db.CheapestRoute(CNF, strings.ToLower(BSB)) // hit
	if db.Cache.Hits != 1 || db.Cache.Misses != 2 {
		t.Errorf(
			"Cache is not working properly: Hits %v Misses %v",
			db.Cache.Hits,
			db.Cache.Misses,
		)
	}
	db.Add(BSB, GRU, 200) // clean
	if db.Cache.cheapestRoute != nil {
		t.Error("Cache is not being cleaned!")
	}
	db.CheapestRoute(APQ, BSB) // don't crash
}
