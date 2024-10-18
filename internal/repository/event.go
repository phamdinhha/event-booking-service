package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/phamdinhha/event-booking-service/internal/model"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
)

type EventRepository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewEventRepository(db *sqlx.DB, logger logger.Logger) EventRepositoryInterface {
	return &EventRepository{db: db, logger: logger}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *model.Event) error {
	query := `
		INSERT INTO events (id, title, description, start_time, end_time, location, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error("EVENT_REPOSITORY.CREATE_EVENT.Error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.Location,
		event.Capacity,
		event.CreatedAt,
		event.UpdatedAt,
	)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error("EVENT_REPOSITORY.CREATE_EVENT.Error", rbErr)
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return fmt.Errorf("failed to create event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("EVENT_REPOSITORY.CREATE_EVENT.Error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, location, capacity, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event model.Event
	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		r.logger.Error("EVENT_REPOSITORY.GET_EVENT_BY_ID.Error", err)
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event *model.Event) error {
	query := `
		UPDATE events
		SET title = $2, description = $3, start_time = $4, end_time = $5, location = $6, capacity = $7, updated_at = $8
		WHERE id = $1
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error("EVENT_REPOSITORY.UPDATE_EVENT.Error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.Location,
		event.Capacity,
		event.UpdatedAt,
	)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error("EVENT_REPOSITORY.UPDATE_EVENT.Error", rbErr)
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		r.logger.Error("EVENT_REPOSITORY.UPDATE_EVENT.Error", err)
		return fmt.Errorf("failed to update event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("EVENT_REPOSITORY.UPDATE_EVENT.Error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error("EVENT_REPOSITORY.DELETE_EVENT.Error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error("EVENT_REPOSITORY.DELETE_EVENT.Error", rbErr)
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		r.logger.Error("EVENT_REPOSITORY.DELETE_EVENT.Error", err)
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("EVENT_REPOSITORY.DELETE_EVENT.Error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
