package controllers

import (
	"bokpang/database"
	"bokpang/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllTickets(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT * FROM tickets")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch tickets",
		})
	}
	defer rows.Close()

	tickets := []models.Ticket{}
	for rows.Next() {
		var ticket models.Ticket
		if err := rows.Scan(&ticket.ID, &ticket.Name, &ticket.Price, &ticket.AvailableSeats, &ticket.CreatedAt, &ticket.UpdatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to parse ticket",
			})
		}
		tickets = append(tickets, ticket)
	}

	return c.JSON(tickets)
}
