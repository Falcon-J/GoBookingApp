package handlers

import (
	"net/http"
	"time"

	"booking-system/database"

	"github.com/gin-gonic/gin"
)

// min returns the minimum of two integers (helper function)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BookingApp holds the database instance and provides HTTP handlers
type BookingApp struct {
	db *database.Database
}

// NewBookingApp creates a new booking application with database
func NewBookingApp() *BookingApp {
	return &BookingApp{
		db: database.NewDatabase(),
	}
}

// HealthCheck returns the health status of the API
func (app *BookingApp) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now(),
	})
}

// GetConferences returns all available conferences
func (app *BookingApp) GetConferences(c *gin.Context) {
	conferences := app.db.GetAllConferences()
	stats := app.db.GetConferenceStats()
	c.JSON(http.StatusOK, gin.H{
		"conferences": conferences,
		"count":       len(conferences),
		"stats":       stats,
	})
}

// CreateUser creates a new user account
func (app *BookingApp) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// If user exists by email, return 409 with existing user to keep entries unique
	if existing, ok := app.db.GetUserByEmail(req.Email); ok {
		c.JSON(http.StatusConflict, existing)
		return
	}

	user, err := app.db.CreateUser(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUserBookings returns all bookings for a specific user
func (app *BookingApp) GetUserBookings(c *gin.Context) {
	userID := c.Param("userID")
	bookings := app.db.GetUserBookings(userID)
	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
		"count":    len(bookings),
	})
}

// CreateBooking creates a new booking (direct booking without reservation)
func (app *BookingApp) CreateBooking(c *gin.Context) {
	var req struct {
		UserID       string `json:"user_id" binding:"required"`
		ConferenceID string `json:"conference_id" binding:"required"`
		TicketCount  int    `json:"ticket_count" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	booking, err := app.db.CreateBooking(req.UserID, req.ConferenceID, req.TicketCount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, booking)
}

// GetBooking retrieves a booking with full details
func (app *BookingApp) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	
	booking := app.db.GetBooking(bookingID)
	if booking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}
	
	// Get additional details
	user, _ := app.db.GetUser(booking.UserID)
	conference, _ := app.db.GetConference(booking.ConferenceID)
	
	c.JSON(http.StatusOK, gin.H{
		"booking":    booking,
		"user":       user,
		"conference": conference,
	})
}

// GetAllBookings returns all bookings with user and conference details
func (app *BookingApp) GetAllBookings(c *gin.Context) {
	bookings := app.db.GetAllBookings()
	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
		"count":    len(bookings),
	})
}



// CreateReservation creates a temporary seat reservation
func (app *BookingApp) CreateReservation(c *gin.Context) {
	var req struct {
		UserID       string `json:"user_id" binding:"required"`
		ConferenceID string `json:"conference_id" binding:"required"`
		TicketCount  int    `json:"ticket_count" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	reservation, err := app.db.CreateReservation(req.UserID, req.ConferenceID, req.TicketCount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	conf, _ := app.db.GetConference(req.ConferenceID)
	c.JSON(http.StatusCreated, gin.H{
		"status":      "success",
		"reservation": reservation,
		"conference":  conf,
		"message":     "Seats reserved for 15 seconds. Complete payment to confirm booking.",
	})
}

// ConfirmReservation converts a reservation to a confirmed booking
func (app *BookingApp) ConfirmReservation(c *gin.Context) {
	reservationID := c.Param("id")
	
	booking, err := app.db.ConfirmReservation(reservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	conf, _ := app.db.GetConference(booking.ConferenceID)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"booking": booking,
		"conference": conf,
		"message": "Payment confirmed! Booking created successfully.",
	})
}

// CancelReservation cancels a seat reservation
func (app *BookingApp) CancelReservation(c *gin.Context) {
	reservationID := c.Param("id")
	
	err := app.db.CancelReservation(reservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Reservation cancelled successfully.",
	})
}

// GetReservation gets a reservation with remaining time
func (app *BookingApp) GetReservation(c *gin.Context) {
	reservationID := c.Param("id")
	
	reservation, err := app.db.GetReservation(reservationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": err.Error()})
		return
	}
	
	// Calculate remaining time
	remainingTime := time.Until(reservation.ExpiresAt)
	if remainingTime < 0 {
		remainingTime = 0
	}
	
	conf, _ := app.db.GetConference(reservation.ConferenceID)
	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"reservation":    reservation,
		"conference":     conf,
		"remaining_time": remainingTime.Seconds(),
		"expired":        remainingTime <= 0,
	})
}

// GetUserReservations gets all active reservations for a user
func (app *BookingApp) GetUserReservations(c *gin.Context) {
	userID := c.Param("userID")
	
	reservations := app.db.GetUserReservations(userID)
	
	// Add remaining time for each reservation
	var result []gin.H
	for _, reservation := range reservations {
		remainingTime := time.Until(reservation.ExpiresAt)
		if remainingTime < 0 {
			remainingTime = 0
		}

		conf, _ := app.db.GetConference(reservation.ConferenceID)
		result = append(result, gin.H{
			"reservation":    reservation,
			"conference":     conf,
			"remaining_time": remainingTime.Seconds(),
			"expired":        remainingTime <= 0,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"reservations": result,
		"count":        len(result),
	})
}

// Queue endpoints
// Enqueue user for conference waitlist
func (app *BookingApp) EnqueueWait(c *gin.Context) {
	var req struct {
		UserID       string `json:"user_id" binding:"required"`
		ConferenceID string `json:"conference_id" binding:"required"`
		TicketCount  int    `json:"ticket_count" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	pos := app.db.EnqueueWait(req.UserID, req.ConferenceID, req.TicketCount)
	c.JSON(http.StatusOK, gin.H{"status": "success", "position": pos})
}

// Get user's queue position
func (app *BookingApp) GetQueuePosition(c *gin.Context) {
	userID := c.Query("user_id")
	conferenceID := c.Param("conferenceID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "user_id required"})
		return
	}
	pos := app.db.GetQueuePosition(userID, conferenceID)
	c.JSON(http.StatusOK, gin.H{"status": "success", "position": pos})
}

// Claim next in queue to create a reservation when it's user's turn
func (app *BookingApp) ClaimNext(c *gin.Context) {
	var req struct {
		UserID       string `json:"user_id" binding:"required"`
		ConferenceID string `json:"conference_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	reservation, err := app.db.ClaimNext(req.UserID, req.ConferenceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	conf, _ := app.db.GetConference(req.ConferenceID)
	c.JSON(http.StatusOK, gin.H{"status": "success", "reservation": reservation, "conference": conf})
}