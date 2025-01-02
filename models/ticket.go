package models

type Ticket struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	AvailableSeats int     `json:"available_seats"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
