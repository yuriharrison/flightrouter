package api

import (
	"github.com/gofiber/fiber"

	"github.com/yuriharrison/flightrouter/flightsdb"
)

var db *flightsdb.FlightsDB

// StartServer start API server on [port]
func StartServer(flightsDB *flightsdb.FlightsDB, port int) {
	db = flightsDB
	app := fiber.New()
	app.Settings.ErrorHandler = ErrorHandler
	AddMiddleware(app)
	AddRoutes(app)
	app.Listen(port)
}
