package api

import (
	"log"

	"github.com/gofiber/fiber"
)

// ErrorHandler API error handler
func ErrorHandler(ctx *fiber.Ctx, err error) {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	log.Println("Error:", err.Error())
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	ctx.Status(code).
		JSON(
			struct {
				Error string `json:"error"`
			}{err.Error()},
		)
}
