package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/google/uuid"
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
	pool := goredis.NewPool(redis)
	return &TicketService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		redis:       redis,
		redsync:     redsync.New(pool),
	}
}

func (s *TicketService) HoldTickets(ctx context.Context, eventID, userID uuid.UUID, quantity int) error {
	mutexName := fmt.Sprintf("lock:event:%d", eventID)
	mutex := s.redsync.NewMutex(
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
	if availableTickets.AvailableTickets < quantity {
		return errors.New("not enough tickets available")
	}

	holdKey := fmt.Sprintf("hold:event:%s:user:%s", eventID.String(), userID.String())
	availableKey := fmt.Sprintf("available:event:%s", eventID.String())

	pipe := s.redis.Pipeline()
	pipe.Set(ctx, holdKey, quantity, bookingHoldTime)
	pipe.DecrBy(ctx, availableKey, int64(quantity))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to hold tickets: %w", err)
	}
	return nil
}

func (s *TicketService) ReleaseHold(ctx context.Context, eventID, userID uuid.UUID) error {
	holdKey := fmt.Sprintf("hold:event:%s:user:%s", eventID.String(), userID.String())
	heldTickets, err := s.redis.Get(ctx, holdKey).Int()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to get held tickets: %w", err)
	}
	if heldTickets > 0 {
		pipe := s.redis.Pipeline()
		pipe.Del(ctx, holdKey)
		pipe.IncrBy(ctx, fmt.Sprintf("available:event:%d", eventID), int64(heldTickets))
		_, err = pipe.Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to release hold: %w", err)
		}
	}
	return nil
}

func (s *TicketService) CleanupExpiredHolds(ctx context.Context) error {
	iter := s.redis.Scan(ctx, 0, "hold:event:*", 0).Iterator()
	for iter.Next(ctx) {
		holdKey := iter.Val()
		ttl := s.redis.TTL(ctx, holdKey).Val()
		if ttl < 0 {
			parts := strings.Split(holdKey, ":")
			if len(parts) == 4 {
				eventID, _ := uuid.Parse(parts[2])
				userID, _ := uuid.Parse(parts[3])
				s.ReleaseHold(ctx, uuid.UUID(eventID), uuid.UUID(userID))
			}
		}
	}
	return nil
}
