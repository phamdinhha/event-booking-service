package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/phamdinhha/event-booking-service/internal/model"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
)

type BookingRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewBookingRepository(db *sqlx.DB, logger logger.Logger) BookingRepositoryInterface {
	return &BookingRepository{
		db:     db,
		logger: logger,
	}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, booking *model.Booking) error {
	query := `
		INSERT INTO bookings (id, event_id, user_id, status, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query,
		booking.ID,
		booking.EventID,
		booking.UserID,
		booking.Status,
		booking.Quantity,
		booking.CreatedAt,
		booking.UpdatedAt,
	)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return fmt.Errorf("failed to create booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *BookingRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*model.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, quantity, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	var booking model.Booking
	err := r.db.GetContext(ctx, &booking, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("booking not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

func (r *BookingRepository) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	query := `
		UPDATE bookings
		SET event_id = $2, user_id = $3, status = $4, quantity = $5, updated_at = $6
		WHERE id = $1
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query,
		booking.ID,
		booking.EventID,
		booking.UserID,
		booking.Status,
		booking.Quantity,
		time.Now(),
	)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return fmt.Errorf("failed to update booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *BookingRepository) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM bookings WHERE id = $1`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *BookingRepository) ListBookings(ctx context.Context, limit, offset int) ([]*model.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, quantity, created_at, updated_at
		FROM bookings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var bookings []*model.Booking
	err := r.db.SelectContext(ctx, &bookings, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}

	return bookings, nil
}
