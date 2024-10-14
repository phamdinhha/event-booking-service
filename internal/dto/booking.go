package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingDTO struct {
	EventID  uuid.UUID `json:"event_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required"`
}

type BookingDTO struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	EventID   uuid.UUID `json:"event_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
	Status    string    `json:"status" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
