package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/phamdinhha/event-booking-service/internal/dto"
	"github.com/phamdinhha/event-booking-service/internal/model"
)

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, booking model.Booking) error
	GetBooking(ctx context.Context, id uuid.UUID) (model.Booking, error)
}

type EventServiceInterface interface {
	CreateEvent(ctx context.Context, eventDTO *dto.CreateEventDTO) (*dto.EventDTO, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*dto.EventDTO, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}
