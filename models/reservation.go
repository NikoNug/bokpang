package models

// Reservation represents a ticket reservation
type Reservation struct {
	ID        int64  `json:"id"`
	TicketID  int64  `json:"ticket_id"`
	UserID    int64  `json:"user_id"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"created_at"`
}
