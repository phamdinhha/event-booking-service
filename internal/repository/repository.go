package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/phamdinhha/event-booking-service/internal/model"
)

type BookingRepositoryInterface interface {
	CreateBooking(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	DeleteBooking(ctx context.Context, id uuid.UUID) error
	ListBookings(ctx context.Context, limit, offset int) ([]*model.Booking, error)
}
