package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

const (
	bookingHoldTime = 5 * time.Minute
)

type EventBookingManager struct {
	redisClient *redis.Client
	redsync     *redsync.Redsync
}

func NewEventBookingManager(redisClient *redis.Client) *EventBookingManager {
	pool := goredis.NewPool(redisClient)
	return &EventBookingManager{
		redisClient: redisClient,
		redsync:     redsync.New(pool),
	}
}

func (m *EventBookingManager) HoldTickets(ctx context.Context, eventID, userID int64, quantity int) error {
	// Create a new mutex for this event
	mutexName := fmt.Sprintf("lock:event:%d", eventID)
	mutex := m.redsync.NewMutex(mutexName,
		redsync.WithExpiry(10*time.Second),
		redsync.WithTries(5))

	// Try to obtain the lock
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer mutex.Unlock()

	// Check available tickets
	availableTickets, err := m.getAvailableTickets(ctx, eventID)
	if err != nil {
		return err
	}

	if availableTickets < quantity {
		return errors.New("not enough tickets available")
	}

	// Hold the tickets
	holdKey := fmt.Sprintf("hold:event:%d:user:%d", eventID, userID)
	availableKey := fmt.Sprintf("available:event:%d", eventID)

	pipe := m.redisClient.Pipeline()
	pipe.Set(ctx, holdKey, quantity, bookingHoldTime)
	pipe.DecrBy(ctx, availableKey, int64(quantity))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to hold tickets: %w", err)
	}

	return nil
}

func (m *EventBookingManager) getAvailableTickets(ctx context.Context, eventID int64) (int, error) {
	availableKey := fmt.Sprintf("available:event:%d", eventID)
	tickets, err := m.redisClient.Get(ctx, availableKey).Int()
	if err == redis.Nil {
		// Key doesn't exist, you might want to initialize it from your database
		return 0, fmt.Errorf("event %d not found", eventID)
	} else if err != nil {
		return 0, fmt.Errorf("failed to get available tickets: %w", err)
	}
	return tickets, nil
}
