package controllers

import (
	"bokpang/database"
	"bokpang/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CreateReservation(c *fiber.Ctx) error {
	reservation := new(models.Reservation)
	if err := c.BodyParser(reservation); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate input
	if reservation.TicketID <= 0 || reservation.UserID <= 0 || reservation.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Ticket ID, User ID, and Quantity are required",
		})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	// Check if the user exists
	var userExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", reservation.UserID).Scan(&userExists)
	if err != nil || !userExists {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Check ticket availability
	var availableSeats int
	err = tx.QueryRow("SELECT available_seats FROM tickets WHERE id = ? FOR UPDATE", reservation.TicketID).Scan(&availableSeats)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to retrieve ticket: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
	}

	if availableSeats < reservation.Quantity {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Not enough available seats"})
	}

	// Reduce available seats
	_, err = tx.Exec("UPDATE tickets SET available_seats = available_seats - ? WHERE id = ?", reservation.Quantity, reservation.TicketID)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to update available seats: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update ticket availability"})
	}

	// Insert the reservation
	query := "INSERT INTO reservations (ticket_id, user_id, quantity) VALUES (?, ?, ?)"
	result, err := tx.Exec(query, reservation.TicketID, reservation.UserID, reservation.Quantity)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to insert reservation: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create reservation"})
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}

	// Return the reservation ID
	lastInsertID, _ := result.LastInsertId()
	return c.Status(201).JSON(fiber.Map{
		"id":        lastInsertID,
		"ticket_id": reservation.TicketID,
		"user_id":   reservation.UserID,
		"quantity":  reservation.Quantity,
		"message":   "Reservation successful",
	})
}
