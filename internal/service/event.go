package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/phamdinhha/event-booking-service/internal/dto"
	"github.com/phamdinhha/event-booking-service/internal/model"
	"github.com/phamdinhha/event-booking-service/internal/repository"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type EventService struct {
	eventRepo repository.EventRepositoryInterface
	logger    logger.Logger
	redis     *redis.Client
}

func NewEventService(
	eventRepo repository.EventRepositoryInterface,
	logger logger.Logger,
	redis *redis.Client,
) *EventServiceInterface {
	return &EventService{
		eventRepo: eventRepo,
		logger:    logger,
		redis:     redis,
	}
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	eventDTO *dto.CreateEventDTO,
) (*dto.EventDTO, error) {
	event := &model.Event{
		ID:               uuid.New(),
		Title:            eventDTO.Title,
		Description:      eventDTO.Description,
		StartTime:        eventDTO.StartTime,
		EndTime:          eventDTO.EndTime,
		Location:         eventDTO.Location,
		Capacity:         eventDTO.Capacity,
		Price:            eventDTO.Price,
		OrganizerId:      eventDTO.OrganizerId,
		CategoryId:       eventDTO.CategoryId,
		Status:           eventDTO.Status,
		AvailableTickets: eventDTO.AvailableTickets,
		CreatedAt:        time.Now(),
	}

	if err := s.eventRepo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}

	return &dto.EventDTO{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Location:    event.Location,
		Capacity:    event.Capacity,
		Price:       event.Price,
		OrganizerId: event.OrganizerId,
		CategoryId:  event.CategoryId,
		Status:      event.Status,
	}, nil
}

func (s *EventService) GetEventByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.EventDTO, error) {
	event, err := s.eventRepo.GetEventByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.EventDTO{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Location:    event.Location,
		Capacity:    event.Capacity,
	}, nil
}

func (s *EventService) DeleteEvent(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.eventRepo.DeleteEvent(ctx, id)
}
