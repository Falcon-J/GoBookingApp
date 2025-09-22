package models

import "time"

// User represents a user in the booking system
type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

// Conference represents a conference that can be booked
type Conference struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Location         string    `json:"location"`
	TotalTickets     int       `json:"total_tickets"`
	AvailableTickets int       `json:"available_tickets"`
	Price            float64   `json:"price"`
	Date             time.Time `json:"date"`
}

// Booking represents a booking made by a user for a conference
type Booking struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ConferenceID  string    `json:"conference_id"`
	TicketsBooked int       `json:"tickets_booked"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`
	BookedAt      time.Time `json:"booked_at"`
}

// SeatReservation represents a temporary seat hold during payment
type SeatReservation struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ConferenceID string    `json:"conference_id"`
	TicketCount  int       `json:"ticket_count"`
	TotalAmount  float64   `json:"total_amount"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

