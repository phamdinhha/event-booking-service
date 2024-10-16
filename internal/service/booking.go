package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phamdinhha/event-booking-service/internal/dto"
	"github.com/phamdinhha/event-booking-service/internal/model"
	"github.com/phamdinhha/event-booking-service/internal/repository"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type BookingService struct {
	bookingRepo repository.BookingRepositoryInterface
	logger      logger.Logger
	redis       *redis.Client
}

func NewBookingService(
	bookingRepo repository.BookingRepositoryInterface,
	logger logger.Logger,
	redis *redis.Client,
) BookingServiceInterface {
	return &BookingService{
		bookingRepo: bookingRepo,
		logger:      logger,
		redis:       redis,
	}
}

func (s *BookingService) CreateBooking(
	ctx context.Context,
	bookDTO *dto.CreateBookingDTO,
) (*dto.BookingDTO, error) {
	holdKey := fmt.Sprintf("hold:event:%s:user:%s", bookDTO.EventID.String(), bookDTO.UserID.String())
	heldTickets, err := s.redis.Get(ctx, holdKey).Int()
	if err != nil {
		return nil, fmt.Errorf("failed to get held tickets: %w", err)
	}

	booking := &model.Booking{
		EventID:   bookDTO.EventID,
		UserID:    bookDTO.UserID,
		Quantity:  heldTickets,
		Status:    "confirmed",
		CreatedAt: time.Now(),
	}

	createdBooking, err := s.bookingRepo.CreateBooking(ctx, booking)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}
	s.cacheBooking(ctx, booking)
	// Release the hold
	s.redis.Del(ctx, holdKey)
	return &dto.BookingDTO{
		ID:        createdBooking.ID,
		EventID:   createdBooking.EventID,
		UserID:    createdBooking.UserID,
		Quantity:  createdBooking.Quantity,
		Status:    createdBooking.Status,
		CreatedAt: createdBooking.CreatedAt,
		UpdatedAt: createdBooking.UpdatedAt,
	}, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*dto.BookingDTO, error) {
	// Try to get from cache first
	cachedBooking, err := s.getCachedBooking(ctx, id)
	if err == nil {
		return &dto.BookingDTO{
			ID:        cachedBooking.ID,
			EventID:   cachedBooking.EventID,
			UserID:    cachedBooking.UserID,
			Quantity:  cachedBooking.Quantity,
			Status:    cachedBooking.Status,
			CreatedAt: cachedBooking.CreatedAt,
			UpdatedAt: cachedBooking.UpdatedAt,
		}, nil
	}
	// If not in cache, get from database
	booking, err := s.bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	// Cache the booking for future requests
	s.cacheBooking(ctx, booking)

	return &dto.BookingDTO{
		ID:        booking.ID,
		EventID:   booking.EventID,
		UserID:    booking.UserID,
		Quantity:  booking.Quantity,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	s.cacheBooking(ctx, booking)
	return nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	if err := s.bookingRepo.DeleteBooking(ctx, id); err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}
	s.invalidateCache(ctx, id)
	return nil
}

func (s *BookingService) ListBookings(ctx context.Context, limit, offset int) ([]*model.Booking, error) {
	bookings, err := s.bookingRepo.ListBookings(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}
	return bookings, nil
}

func (s *BookingService) cacheBooking(ctx context.Context, booking *model.Booking) {
	bookingJSON, err := json.Marshal(booking)
	if err != nil {
		s.logger.Error("failed to marshal booking for caching", "error", err)
		return
	}
	err = s.redis.Set(ctx, fmt.Sprintf("booking:%s", booking.ID), bookingJSON, time.Hour).Err()
	if err != nil {
		s.logger.Error("failed to cache booking", "error", err)
	}
}

func (s *BookingService) getCachedBooking(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	bookingJSON, err := s.redis.Get(ctx, fmt.Sprintf("booking:%s", id)).Result()
	if err != nil {
		return nil, err
	}

	var booking model.Booking
	err = json.Unmarshal([]byte(bookingJSON), &booking)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached booking: %w", err)
	}

	return &booking, nil
}

func (s *BookingService) invalidateCache(ctx context.Context, id uuid.UUID) {
	err := s.redis.Del(ctx, fmt.Sprintf("booking:%s", id)).Err()
	if err != nil {
		s.logger.Error("failed to invalidate booking cache", "error", err)
	}
}
