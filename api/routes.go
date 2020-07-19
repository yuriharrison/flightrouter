package api

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

// AddRoutes Adds API routes to *fiber.App
func AddRoutes(app *fiber.App) {
	app.Static("/", "./public")
	app.Get("/ws/status", websocket.New(StatusWSView))

	app.Put("/flight", CreateFlightView)
	app.Delete("/flight/:route", DeleteFlightView)
	app.Put("/flights/import/csv", ImportFlightsView)
	app.Get("/flights/search/:route", CheapestRouteView)

	app.Use(NotFoundView)
}
