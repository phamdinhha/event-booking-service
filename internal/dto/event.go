package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventDTO struct {
	Title            string    `json:"title" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	StartTime        time.Time `json:"start_time" validate:"required"`
	EndTime          time.Time `json:"end_time" validate:"required"`
	Location         string    `json:"location" validate:"required"`
	Capacity         int       `json:"capacity" validate:"required"`
	Price            float64   `json:"price" validate:"required"`
	OrganizerId      uuid.UUID `json:"organizer_id" validate:"required"`
	CategoryId       uuid.UUID `json:"category_id" validate:"required"`
	Status           string    `json:"status" validate:"required"`
	AvailableTickets int       `json:"available_tickets" validate:"required"`
}

type EventDTO struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Location         string    `json:"location"`
	Capacity         int       `json:"capacity"`
	Price            float64   `json:"price"`
	OrganizerId      uuid.UUID `json:"organizer_id"`
	CategoryId       uuid.UUID `json:"category_id"`
	Status           string    `json:"status"`
	AvailableTickets int       `json:"available_tickets"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
