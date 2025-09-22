package database

import (
	"booking-system/models"
	"testing"
)

// helper to build DB with a user and conference
func makeDBWithUserAndConf() (*Database, *models.User, *models.Conference) {
	db := NewDatabase()
	u, _ := db.CreateUser("Alice", "alice@example.com")
	var conf *models.Conference
	for _, c := range db.Conferences {
		if c.ID == "conf-1" {
			conf = c
			break
		}
	}
	if conf == nil {
		// take any
		for _, c := range db.Conferences {
			conf = c
			break
		}
	}
	return db, u, conf
}

func TestUserEmailUniqueness(t *testing.T) {
	db := NewDatabase()
	_, err := db.CreateUser("A", "Test@Example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = db.CreateUser("B", "test@example.com")
	if err == nil {
		t.Fatalf("expected duplicate email error, got nil")
	}
}

func TestSingleActiveReservationPerUserPerConference(t *testing.T) {
	db, user, conf := makeDBWithUserAndConf()
	// first reservation should succeed
	res1, err := db.CreateReservation(user.ID, conf.ID, 1)
	if err != nil || res1 == nil {
		t.Fatalf("expected first reservation ok, got err=%v", err)
	}
	// second reservation for same conference should fail while first alive
	if _, err := db.CreateReservation(user.ID, conf.ID, 1); err == nil {
		t.Fatalf("expected error for duplicate active reservation")
	}
}
