package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/phamdinhha/event-booking-service/internal/model"
)

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, booking model.Booking) error
	GetBooking(ctx context.Context, id uuid.UUID) (model.Booking, error)
}
