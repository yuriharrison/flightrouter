package api

import (
	"github.com/gofiber/fiber"

	"github.com/yuriharrison/flightrouter/loader"
	"github.com/yuriharrison/flightrouter/util"
)

// NotFoundView - 404
func NotFoundView(c *fiber.Ctx) {
	c.Next(fiber.NewError(fiber.StatusNotFound, "Not found"))
}

// CheapestRouteView calculate the cheapest route for a given origin-destination
func CheapestRouteView(c *fiber.Ctx) {
	route := c.Params("route")
	origin, destination := util.FormatInputRoute(route)
	cheapestRoute, err := db.CheapestRoute(origin, destination)
	if err != nil {
		c.Next(err)
		return
	}
	var totalCost float64
	for _, flight := range cheapestRoute {
		totalCost += flight.Price
	}
	description := util.FormatRouteToString(cheapestRoute)
	c.Status(fiber.StatusOK).
		JSON(
			CheapestRoutePayload{
				TotalCost:   totalCost,
				Description: description,
				Flights:     cheapestRoute,
			},
		)
}

// CreateFlightView Create a new flight
func CreateFlightView(c *fiber.Ctx) {
	flight := &FlightPayload{}
	if err := c.BodyParser(flight); err != nil {
		c.Next(err)
		return
	}
	err := db.Add(flight.Origin, flight.Destination, flight.Price)
	if err != nil {
		c.Next(err)
		return
	}
	c.Status(fiber.StatusCreated).JSON(flight)
}

// DeleteFlightView Delete a given flight
func DeleteFlightView(c *fiber.Ctx) {
	route := c.Params("route")
	origin, destination := util.FormatInputRoute(route)
	if err := db.Remove(origin, destination); err != nil {
		c.Next(err)
		return
	}
	c.Status(fiber.StatusNoContent)
}

// ImportFlightsView Import a CSV file
func ImportFlightsView(c *fiber.Ctx) {
	formFile, err := c.FormFile("document")
	if err != nil {
		c.Next(err)
		return
	}
	file, err := formFile.Open()
	if err != nil {
		c.Next(err)
		return
	}
	loader.LoadFlights(file, db)
	c.Status(fiber.StatusNoContent)
}
