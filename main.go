package main

import (
	"bokpang/database"
	"bokpang/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()

	routes.TicketRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
