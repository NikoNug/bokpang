package controllers

import (
	"bokpang/database"
	"bokpang/models"
	"log"

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

func CreateTicket(c *fiber.Ctx) error {
	ticket := new(models.Ticket)
	if err := c.BodyParser(ticket); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse json",
		})
	}

	if ticket.Name == "" || ticket.Price <= 0 || ticket.AvailableSeats <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name, Price, and Available Seats are required",
		})
	}

	result, err := database.DB.Exec("INSERT INTO tickets (name, price, available_seats) VALUES (?,?,?)", ticket.Name, ticket.Price, ticket.AvailableSeats)
	if err != nil {
		log.Printf("Failed to insert ticket : %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create ticket",
		})
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to retreive inserted ID : %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retreive inserted ID",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":              lastInsertID,
		"name":            ticket.Name,
		"price":           ticket.Price,
		"available_seats": ticket.AvailableSeats,
	})
}
