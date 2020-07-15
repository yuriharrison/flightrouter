package flightsdb

import (
	"fmt"
	"strings"
	"testing"
)

var codes []string = []string{"APQ", "BAZ", "BEL", "PLU", "BSB", "VCP", "CPQ", "CAU", "CAF", "CAC", "CIZ", "CDJ", "CZS", "BFH", "CWB", "FEJ", "FLN", "FOR", "IGU", "IZA", "GYN", "GRU", "IMP"}
var dataSize int = 10
var data = []struct {
	orig, dest string
	price      float32
}{
	{"cnf", "gru", 200},
	{"cnf", "gru", 215.5},
	{"cnf", "apq", 899.99},
	{"gru", "cnf", 70},
	{"gru", "cnf", 78},
	{"gru", "baz", 300},
	{"gru", "plu", 10},
	{"baz", "bel", 600},
	{"baz", "gru", 10},
	{"baz", "plu", 200},
	{"plu", "apq", 15},
	{"bsb", "vcp", 415.5},
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
			routeCodes = append(routeCodes, flight.orig.Code)
			if !fail && flight.orig.Code != strings.ToUpper(p.expectedRoute[i]) {
				fail = true
			}
		}
		routeCodes = append(routeCodes, flight.dest.Code)
		if !fail && flight.dest.Code != strings.ToUpper(p.expectedRoute[i+1]) {
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
	db := NewFlightsDB()
	// Assert size
	for _, d := range data {
		db.Add(d.orig, d.dest, d.price)
	}
	if l := db.Size(); l != dataSize {
		t.Errorf("FlightsDB.Size() %v != %v expected", l, 1)
	}

	// Assert cheapest route
	assertRoute(assertRoutePayload{
		t, db, "baz", "bel", []string{"BAZ", "BEL"},
	})
	assertRoute(assertRoutePayload{
		t, db, "cnf", "apq", []string{"CNF", "GRU", "PLU", "APQ"},
	})
	db.Remove("GRU", "PLU")
	assertRoute(assertRoutePayload{
		t, db, "cnf", "apq", []string{"CNF", "GRU", "BAZ", "PLU", "APQ"},
	})
	db.Add("baz", "plu", 400)
	assertRoute(assertRoutePayload{
		t, db, "cnf", "apq", []string{"CNF", "APQ"},
	})

	// Assert no route available
	if route, err := db.CheapestRoute("CNF", "BSB"); err == nil || len(route) > 0 {
		t.Error(err)
	}
	if route, err := db.CheapestRoute("AAA", "ZZZ"); err == nil || len(route) > 0 {
		t.Error(err)
	}

	// Remove invalid
	db.Remove("AAA", "ZZZ")
}
