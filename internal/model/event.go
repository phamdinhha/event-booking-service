package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Title            string    `json:"title" db:"title"`
	Description      string    `json:"description" db:"description"`
	StartTime        time.Time `json:"start_time" db:"start_time"`
	EndTime          time.Time `json:"end_time" db:"end_time"`
	Location         string    `json:"location" db:"location"`
	Capacity         int       `json:"capacity" db:"capacity"`
	AvailableTickets int       `json:"available_tickets" db:"available_tickets"`
	Price            float64   `json:"price" db:"price"`
	OrganizerId      uuid.UUID `json:"organizer_id" db:"organizer_id"`
	CategoryId       uuid.UUID `json:"category_id" db:"category_id"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
)

func (s EventStatus) Validate() bool {
	switch s {
	case EventStatusDraft, EventStatusPublished, EventStatusCancelled, EventStatusCompleted:
		return true
	}
	return false
}

// CheckTicketAvailability verifies if the requested number of tickets is available
func (e *Event) CheckTicketAvailability(requestedTickets int) error {
	if requestedTickets <= 0 {
		return errors.New("requested tickets must be greater than zero")
	}
	if requestedTickets > e.AvailableTickets {
		return errors.New("not enough tickets available")
	}
	return nil
}

// ReserveTickets attempts to reserve the specified number of tickets
func (e *Event) ReserveTickets(tickets int) error {
	if err := e.CheckTicketAvailability(tickets); err != nil {
		return err
	}
	e.AvailableTickets -= tickets
	return nil
}

// ReleaseTickets releases the specified number of tickets back to the available pool
func (e *Event) ReleaseTickets(tickets int) error {
	if tickets <= 0 {
		return errors.New("number of tickets to release must be greater than zero")
	}
	if e.AvailableTickets+tickets > e.Capacity {
		return errors.New("cannot release more tickets than the event capacity")
	}
	e.AvailableTickets += tickets
	return nil
}
