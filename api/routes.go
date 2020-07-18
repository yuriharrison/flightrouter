package api

import (
	"github.com/gofiber/fiber"
)

// AddRoutes Adds API routes to *fiber.App
func AddRoutes(app *fiber.App) {
	app.Put("/flight", CreateFlightView)
	app.Delete("/flight/:route", DeleteFlightView)
	app.Put("/flights/import/csv", ImportFlightsView)
	app.Get("/flights/search/:route", CheapestRouteView)
	app.Use(NotFoundView)
}
