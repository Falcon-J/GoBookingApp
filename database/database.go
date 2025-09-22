package database

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"booking-system/models"

	"github.com/google/uuid"
)

// Database represents an in-memory database for the booking system
type Database struct {
	Users         map[string]*models.User
	Conferences   map[string]*models.Conference
	Bookings      map[string]*models.Booking
	Reservations  map[string]*models.SeatReservation
	WaitQueues    map[string][]*WaitEntry // per-conference wait queues
	StartTime     time.Time        // Track when the database was initialized
	mutex         sync.RWMutex     // Thread-safe operations
}

// WaitEntry represents a queued request for tickets
type WaitEntry struct {
	ID           string
	UserID       string
	ConferenceID string
	TicketCount  int
	EnqueuedAt   time.Time
}

// NewDatabase creates a new database instance with sample data
func NewDatabase() *Database {
	db := &Database{
		Users:         make(map[string]*models.User),
		Conferences:   make(map[string]*models.Conference),
		Bookings:      make(map[string]*models.Booking),
		Reservations:  make(map[string]*models.SeatReservation),
		WaitQueues:    make(map[string][]*WaitEntry),
		StartTime:     time.Now(),
	}
	
	// Add sample data
	db.addSampleData()
	return db
}

// addSampleData populates the database with sample conferences
func (db *Database) addSampleData() {
	// Add sample conferences
	conf1 := &models.Conference{
		ID:               "conf-1",
		Name:             "Go Conference 2024",
		Location:         "San Francisco",
		TotalTickets:     100,
		AvailableTickets: 100,
		Price:            299.99,
		Date:             time.Now().AddDate(0, 2, 0), // 2 months from now
	}
	
	conf2 := &models.Conference{
		ID:               "conf-2",
		Name:             "DevOps Summit",
		Location:         "New York",
		TotalTickets:     75,
		AvailableTickets: 75,
		Price:            399.99,
		Date:             time.Now().AddDate(0, 3, 0), // 3 months from now
	}
	
	conf3 := &models.Conference{
		ID:               "conf-3",
		Name:             "Cloud Native Expo",
		Location:         "Seattle",
		TotalTickets:     150,
		AvailableTickets: 150,
		Price:            199.99,
		Date:             time.Now().AddDate(0, 1, 15), // 1.5 months from now
	}
	
	db.Conferences[conf1.ID] = conf1
	db.Conferences[conf2.ID] = conf2
	db.Conferences[conf3.ID] = conf3
	
	log.Printf("Added %d sample conferences to database", len(db.Conferences))
}

// CreateUser creates a new user in the database
func (db *Database) CreateUser(name, email string) (*models.User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Normalize email for uniqueness (case-insensitive)
	norm := strings.ToLower(strings.TrimSpace(email))

	// Check if user with email already exists
	for _, user := range db.Users {
		if strings.ToLower(strings.TrimSpace(user.Email)) == norm {
			return nil, fmt.Errorf("user with email %s already exists", email)
		}
	}

	user := &models.User{
		ID:      uuid.New().String(),
		Name:    name,
		Email:   norm,
		Created: time.Now(),
	}
	
	db.Users[user.ID] = user
	return user, nil
}

// GetUser retrieves a user by ID
func (db *Database) GetUser(userID string) (*models.User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	user, exists := db.Users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	
	return user, nil
}

// GetAllConferences returns all conferences
func (db *Database) GetAllConferences() []*models.Conference {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	var conferences []*models.Conference
	for _, conf := range db.Conferences {
		conferences = append(conferences, conf)
	}
	// Keep conferences sorted by ID for convenience
	sort.Slice(conferences, func(i, j int) bool {
		return conferences[i].ID < conferences[j].ID
	})
	return conferences
}

// GetConference retrieves a conference by ID
func (db *Database) GetConference(conferenceID string) (*models.Conference, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	conference, exists := db.Conferences[conferenceID]
	if !exists {
		return nil, fmt.Errorf("conference not found")
	}
	
	return conference, nil
}

// CreateBooking creates a new booking
func (db *Database) CreateBooking(userID, conferenceID string, ticketCount int) (*models.Booking, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	conference, exists := db.Conferences[conferenceID]
	if !exists {
		return nil, fmt.Errorf("conference not found")
	}
	
	if conference.AvailableTickets < ticketCount {
		return nil, fmt.Errorf("not enough tickets available")
	}
	
	booking := &models.Booking{
		ID:            uuid.New().String(),
		UserID:        userID,
		ConferenceID:  conferenceID,
		TicketsBooked: ticketCount,
		TotalAmount:   conference.Price * float64(ticketCount),
		Status:        "confirmed",
		BookedAt:      time.Now(),
	}
	
	// Update available tickets
	conference.AvailableTickets -= ticketCount
	
	db.Bookings[booking.ID] = booking
	return booking, nil
}

// GetUserBookings returns all bookings for a user
func (db *Database) GetUserBookings(userID string) []*models.Booking {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	var bookings []*models.Booking
	for _, booking := range db.Bookings {
		if booking.UserID == userID {
			bookings = append(bookings, booking)
		}
	}
	return bookings
}

// GetBooking retrieves a booking by ID
func (db *Database) GetBooking(id string) *models.Booking {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	booking, exists := db.Bookings[id]
	if !exists {
		return nil
	}
	
	return booking
}

// GetAllBookings returns all bookings with user and conference details
func (db *Database) GetAllBookings() []map[string]interface{} {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	var result []map[string]interface{}
	
	for _, booking := range db.Bookings {
		user := db.Users[booking.UserID]
		conference := db.Conferences[booking.ConferenceID]
		
		bookingData := map[string]interface{}{
			"booking":    booking,
			"user":       user,
			"conference": conference,
		}
		result = append(result, bookingData)
	}
	
	return result
}

// GetAllUsers returns all users
func (db *Database) GetAllUsers() []*models.User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	var users []*models.User
	for _, user := range db.Users {
		users = append(users, user)
	}
	return users
}

// ResetDatabase clears all data and reinitializes with sample data
func (db *Database) ResetDatabase() {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	// Clear all maps
	db.Users = make(map[string]*models.User)
	db.Conferences = make(map[string]*models.Conference)
	db.Bookings = make(map[string]*models.Booking)
	db.Reservations = make(map[string]*models.SeatReservation)
	// admin sessions removed
	db.WaitQueues = make(map[string][]*WaitEntry)
	
	// Reset start time
	db.StartTime = time.Now()
	
	// Repopulate with sample data
	db.addSampleData()
}

// CreateReservation creates a temporary seat reservation
func (db *Database) CreateReservation(userID, conferenceID string, ticketCount int) (*models.SeatReservation, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Clean up expired reservations first (already holding write lock)
	db.cleanupExpiredReservationsLocked()
	
	conference, exists := db.Conferences[conferenceID]
	if !exists {
		return nil, fmt.Errorf("conference not found")
	}
	
	// Ensure user has no other active reservation for this conference
	for _, reservation := range db.Reservations {
		if reservation.UserID == userID && reservation.ConferenceID == conferenceID {
			if time.Now().Before(reservation.ExpiresAt) {
				return nil, fmt.Errorf("you already have an active reservation for this conference")
			}
		}
	}

	// Calculate total reserved tickets for this conference
	reservedTickets := 0
	for _, reservation := range db.Reservations {
		if reservation.ConferenceID == conferenceID {
			reservedTickets += reservation.TicketCount
		}
	}
	
	// Check if enough tickets are available (considering reservations)
	availableForReservation := conference.AvailableTickets - reservedTickets
	if availableForReservation < ticketCount {
		return nil, fmt.Errorf("not enough tickets available for reservation")
	}
	
	reservation := &models.SeatReservation{
		ID:           uuid.New().String(),
		UserID:       userID,
		ConferenceID: conferenceID,
		TicketCount:  ticketCount,
		TotalAmount:  conference.Price * float64(ticketCount),
		ExpiresAt:    time.Now().Add(15 * time.Second),
		CreatedAt:    time.Now(),
	}
	
	db.Reservations[reservation.ID] = reservation
	return reservation, nil
}

// ConfirmReservation converts a reservation to a booking
func (db *Database) ConfirmReservation(reservationID string) (*models.Booking, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	reservation, exists := db.Reservations[reservationID]
	if !exists {
		return nil, fmt.Errorf("reservation not found")
	}
	
	// Check if reservation has expired
	if time.Now().After(reservation.ExpiresAt) {
		delete(db.Reservations, reservationID)
		return nil, fmt.Errorf("reservation has expired")
	}
	
	// Create the booking
	booking := &models.Booking{
		ID:            uuid.New().String(),
		UserID:        reservation.UserID,
		ConferenceID:  reservation.ConferenceID,
		TicketsBooked: reservation.TicketCount,
		TotalAmount:   reservation.TotalAmount,
		Status:        "confirmed",
		BookedAt:      time.Now(),
	}
	
	// Update conference availability
	conference := db.Conferences[reservation.ConferenceID]
	conference.AvailableTickets -= reservation.TicketCount
	
	// Store booking and remove reservation
	db.Bookings[booking.ID] = booking
	delete(db.Reservations, reservationID)
	
	return booking, nil
}

// CancelReservation removes a reservation
func (db *Database) CancelReservation(reservationID string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	
	if _, exists := db.Reservations[reservationID]; !exists {
		return fmt.Errorf("reservation not found")
	}
	
	delete(db.Reservations, reservationID)
	return nil
}

// GetReservation gets a reservation by ID
func (db *Database) GetReservation(reservationID string) (*models.SeatReservation, error) {
	// Clean up expired reservations first with exclusive lock
	db.cleanupExpiredReservations()

	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	reservation, exists := db.Reservations[reservationID]
	if !exists {
		return nil, fmt.Errorf("reservation not found")
	}
	
	return reservation, nil
}

// GetUserReservations gets all active reservations for a user
func (db *Database) GetUserReservations(userID string) []*models.SeatReservation {
	// Clean up expired reservations first
	db.cleanupExpiredReservations()

	db.mutex.RLock()
	defer db.mutex.RUnlock()
	
	var reservations []*models.SeatReservation
	for _, reservation := range db.Reservations {
		if reservation.UserID == userID {
			reservations = append(reservations, reservation)
		}
	}
	return reservations
}

// cleanupExpiredReservations removes expired reservations (internal method)
func (db *Database) cleanupExpiredReservations() {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.cleanupExpiredReservationsLocked()
}

// cleanupExpiredReservationsLocked removes expired reservations; caller must hold write lock
func (db *Database) cleanupExpiredReservationsLocked() {
	now := time.Now()
	for id, reservation := range db.Reservations {
		if now.After(reservation.ExpiresAt) {
			delete(db.Reservations, id)
		}
	}
}

// GetUserByEmail returns a user by email (case-insensitive)
func (db *Database) GetUserByEmail(email string) (*models.User, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	norm := strings.ToLower(strings.TrimSpace(email))
	for _, u := range db.Users {
		if strings.ToLower(strings.TrimSpace(u.Email)) == norm {
			return u, true
		}
	}
	return nil, false
}



// GetConferenceStats returns reserved count and queue length per conference
func (db *Database) GetConferenceStats() map[string]struct{ Reserved int; Queue int } {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	// compute reserved counts ignoring expired
	now := time.Now()
	stats := make(map[string]struct{ Reserved int; Queue int })
	for id := range db.Conferences {
		stats[id] = struct{ Reserved int; Queue int }{Reserved: 0, Queue: len(db.WaitQueues[id])}
	}
	for _, r := range db.Reservations {
		if now.Before(r.ExpiresAt) {
			s := stats[r.ConferenceID]
			s.Reserved += r.TicketCount
			stats[r.ConferenceID] = s
		}
	}
	return stats
}

// EnqueueWait adds a user to the conference wait queue, returns 1-based position
func (db *Database) EnqueueWait(userID, conferenceID string, ticketCount int) int {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	q := db.WaitQueues[conferenceID]
	// avoid duplicate entries for same user+conference; keep earliest
	for i, e := range q {
		if e.UserID == userID {
			// update ticketCount to latest request
			q[i].TicketCount = ticketCount
			db.WaitQueues[conferenceID] = q
			return i + 1
		}
	}
	entry := &WaitEntry{
		ID:           uuid.New().String(),
		UserID:       userID,
		ConferenceID: conferenceID,
		TicketCount:  ticketCount,
		EnqueuedAt:   time.Now(),
	}
	q = append(q, entry)
	db.WaitQueues[conferenceID] = q
	return len(q)
}

// GetQueuePosition returns 1-based position, or 0 if not present
func (db *Database) GetQueuePosition(userID, conferenceID string) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	q := db.WaitQueues[conferenceID]
	for i, e := range q {
		if e.UserID == userID {
			return i + 1
		}
	}
	return 0
}

// ClaimNext attempts to create a reservation for the first-in-queue user if they are the caller.
func (db *Database) ClaimNext(userID, conferenceID string) (*models.SeatReservation, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.cleanupExpiredReservationsLocked()
	q := db.WaitQueues[conferenceID]
	if len(q) == 0 || q[0].UserID != userID {
		return nil, fmt.Errorf("not your turn yet")
	}
	conf, ok := db.Conferences[conferenceID]
	if !ok {
		return nil, fmt.Errorf("conference not found")
	}
	// compute currently reserved for this conf
	reserved := 0
	now := time.Now()
	for _, r := range db.Reservations {
		if r.ConferenceID == conferenceID && now.Before(r.ExpiresAt) {
			reserved += r.TicketCount
		}
	}
	available := conf.AvailableTickets - reserved
	need := q[0].TicketCount
	if available < need {
		return nil, fmt.Errorf("not enough tickets available")
	}
	// create reservation
	res := &models.SeatReservation{
		ID:           uuid.New().String(),
		UserID:       userID,
		ConferenceID: conferenceID,
		TicketCount:  need,
		TotalAmount:  conf.Price * float64(need),
		ExpiresAt:    time.Now().Add(15 * time.Second),
		CreatedAt:    time.Now(),
	}
	db.Reservations[res.ID] = res
	// pop queue head
	db.WaitQueues[conferenceID] = q[1:]
	return res, nil
}