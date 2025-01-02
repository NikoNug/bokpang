package routes

import (
	"bokpang/controllers"

	"github.com/gofiber/fiber/v2"
)

func TicketRoutes(app *fiber.App) {
	route := app.Group("/tickets")
	route.Get("/", controllers.GetAllTickets)
}
