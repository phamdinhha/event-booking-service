package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/phamdinhha/event-booking-service/internal/dto"
)

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, booking *dto.CreateBookingDTO) (*dto.BookingDTO, error)
	GetBooking(ctx context.Context, id uuid.UUID) (*dto.BookingDTO, error)
	DeleteBooking(ctx context.Context, id uuid.UUID) error
}

type EventServiceInterface interface {
	CreateEvent(ctx context.Context, eventDTO *dto.CreateEventDTO) (*dto.EventDTO, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*dto.EventDTO, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}
