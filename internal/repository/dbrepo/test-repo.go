package dbrepo

import (
	"errors"
	"time"

	"github.com/trcwrx78/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool  {
	return true
}

// Inserta a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// Inserts a room restiction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// Returns true if availability exists for roomID, and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)  {
	return false, nil
}

// Returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabiltyForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// Get a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error)  {
	var room models.Room
	if id > 2 {
		return room, errors.New("Some error")
	}

	return room, nil
}