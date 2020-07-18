package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	db "github.com/yuriharrison/flightrouter/flightsdb"
)

// ReadJSON Read a json from a standard *io.Reader to an *interface{}
func ReadJSON(reader io.Reader, interf interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	json.Unmarshal(data, interf)
	return nil
}

// FormatInputRoute format  input route e.g. GRU-CNF
func FormatInputRoute(rawRoute string) (string, string) {
	routeList := strings.Split(rawRoute, "-")
	if len(routeList) != 2 {
		fmt.Println("Wrong format. Try [ORIGIN]-[DESTINY] e.g.: GRU-CNF")
		os.Exit(1)
	}
	return routeList[0], routeList[1][:3]
}

// FormatRouteToString format a list of *Flights to string
func FormatRouteToString(route []*db.Flight) string {
	str := ""
	var fullPrice float64
	for i, flight := range route {
		if i == 0 {
			str += flight.Origin.Code + "-" + flight.Destination.Code
		} else {
			str += "-" + flight.Destination.Code
		}
		fullPrice += flight.Price
	}
	return str + fmt.Sprintf(" -> $%.2f", fullPrice)
}
