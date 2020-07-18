package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/gofiber/utils"

	"github.com/yuriharrison/flightrouter/flightsdb"
	"github.com/yuriharrison/flightrouter/loader"
	"github.com/yuriharrison/flightrouter/util"
)

const ContentType = "Content-Type"

func New() (app *fiber.App) {
	app = fiber.New()
	db = flightsdb.New()
	loader.ImportFlightsFromFile("../fixtures/flights.csv", db)
	AddMiddleware(app)
	AddRoutes(app)
	return
}

func payloadReader(payload interface{}) io.Reader {
	data, _ := json.Marshal(payload)
	return bytes.NewReader(data)
}

func assertRequest(t *testing.T,
	app *fiber.App,
	request *http.Request,
	status int,
	expected interface{}) {
	resp, err := app.Test(request)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, status, resp.StatusCode)
	payload := make(map[string]interface{})
	if expected == nil {
		return
	}
	utils.AssertEqual(
		t,
		fiber.MIMEApplicationJSON,
		resp.Header.Get(ContentType),
	)
	util.ReadJSON(resp.Body, &payload)
	expectedAsByte, _ := json.Marshal(expected)
	expectedAsMap := make(map[string]interface{})
	json.Unmarshal(expectedAsByte, &expectedAsMap)
	if !reflect.DeepEqual(payload, expectedAsMap) {
		fmt.Println("Expected:\t", expectedAsMap)
		fmt.Println("Got:\t\t", payload)
		t.Error("Payloads differ!")
	}
}

func TestCheapestRouteView(t *testing.T) {
	app := New()
	assertRequest(
		t,
		app,
		httptest.NewRequest("GET", "/flights/search/GRU-APQ", nil),
		fiber.StatusOK,
		CheapestRoutePayload{
			TotalCost:   25,
			Description: "GRU-PLU-APQ -> $25.00",
			Flights: []*flightsdb.Flight{
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "GRU"},
					Destination: &flightsdb.Airport{Code: "PLU"},
					Price:       10,
				},
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "PLU"},
					Destination: &flightsdb.Airport{Code: "APQ"},
					Price:       15,
				},
			},
		},
	)
}

func TestCreateFlight(t *testing.T) {
	app := New()
	flightPayload := FlightPayload{
		Origin:      "GRU",
		Destination: "APQ",
		Price:       10,
	}
	addFlightRequest := httptest.NewRequest(
		"PUT", "/flight", payloadReader(flightPayload))
	addFlightRequest.Header.Add(ContentType, fiber.MIMEApplicationJSON)
	assertRequest(t, app, addFlightRequest, fiber.StatusCreated, flightPayload)
	assertRequest(
		t,
		app,
		httptest.NewRequest("GET", "/flights/search/GRU-APQ", nil),
		fiber.StatusOK,
		CheapestRoutePayload{
			TotalCost:   10,
			Description: "GRU-APQ -> $10.00",
			Flights: []*flightsdb.Flight{
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "GRU"},
					Destination: &flightsdb.Airport{Code: "APQ"},
					Price:       10,
				},
			},
		},
	)
}

func TestDeleteFlight(t *testing.T) {
	app := New()
	assertRequest(
		t,
		app,
		httptest.NewRequest("DELETE", "/flight/GRU-PLU", nil),
		fiber.StatusNoContent,
		nil,
	)
	assertRequest(
		t,
		app,
		httptest.NewRequest("GET", "/flights/search/GRU-APQ", nil),
		fiber.StatusOK,
		CheapestRoutePayload{
			TotalCost:   515,
			Description: "GRU-BAZ-PLU-APQ -> $515.00",
			Flights: []*flightsdb.Flight{
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "GRU"},
					Destination: &flightsdb.Airport{Code: "BAZ"},
					Price:       300,
				},
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "BAZ"},
					Destination: &flightsdb.Airport{Code: "PLU"},
					Price:       200,
				},
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "PLU"},
					Destination: &flightsdb.Airport{Code: "APQ"},
					Price:       15,
				},
			},
		},
	)
}

func TestImportFlights(t *testing.T) {
	app := New()
	body := new(bytes.Buffer)
	form := multipart.NewWriter(body)
	_, err := form.CreateFormFile("document", "../fixtures/flights.csv")
	if err != nil {
		t.Error(err)
	}
	form.Close()
	request := httptest.NewRequest("PUT", "/flights/import/csv", body)
	request.Header.Add(ContentType, form.FormDataContentType())
	assertRequest(t, app, request, fiber.StatusNoContent, nil)
	assertRequest(
		t,
		app,
		httptest.NewRequest("GET", "/flights/search/GRU-CNF", nil),
		fiber.StatusOK,
		CheapestRoutePayload{
			TotalCost:   78,
			Description: "GRU-CNF -> $78.00",
			Flights: []*flightsdb.Flight{
				&flightsdb.Flight{
					Origin:      &flightsdb.Airport{Code: "GRU"},
					Destination: &flightsdb.Airport{Code: "CNF"},
					Price:       78,
				},
			},
		},
	)
}

func TestErrorView(t *testing.T) {
	app := New()
	assertRequest(
		t,
		app,
		httptest.NewRequest("GET", "/flights/search/GRUA-APQ", nil),
		fiber.StatusInternalServerError,
		nil,
	)
	assertRequest(
		t,
		app,
		httptest.NewRequest("GET", "/notfound", nil),
		fiber.StatusNotFound,
		nil,
	)
}
