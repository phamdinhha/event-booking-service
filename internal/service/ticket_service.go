package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/phamdinhha/event-booking-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

const (
	bookingHoldTime = 5 * time.Minute
)

type TicketService struct {
	bookingRepo repository.BookingRepositoryInterface
	eventRepo   repository.EventRepositoryInterface
	redis       *redis.Client
	redsync     *redsync.Redsync
}

func NewTicketService(
	bookingRepo repository.BookingRepositoryInterface,
	eventRepo repository.EventRepositoryInterface,
	redis *redis.Client,
) *TicketService {
	return &TicketService{bookingRepo, eventRepo, redis}
}

func (s *TicketService) HoldTickets(ctx context.Context, eventID, userID int64, quantity int) error {
	mutexName := fmt.Sprintf("lock:event:%d", eventID)
	mutex := m.redsync.NewMutex(
		mutexName,
		redsync.WithExpiry(10*time.Second),
		redsync.WithTries(5),
	)
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer mutex.Unlock()

	availableTickets, err := s.eventRepo.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}
	if availableTickets.TotalTickets < quantity {
		return errors.New("not enough tickets available")
	}

	holdKey := fmt.Sprintf("hold:event:%d:user:%d", eventID, userID)
	availableKey := fmt.Sprintf("available:event:%d", eventID)
}
