package api

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

// AddMiddleware Adds middleware to *fiber.App
func AddMiddleware(app *fiber.App) {
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
}
