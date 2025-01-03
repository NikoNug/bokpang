package routes

import (
	"bokpang/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(app *fiber.App) {
	ticket := app.Group("/tickets")
	ticket.Get("/", controllers.GetAllTickets)
	ticket.Post("/", controllers.CreateTicket)

	reservation := app.Group("/reservation")
	reservation.Post("/", controllers.CreateReservation)
}
