package http_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/phamdinhha/event-booking-service/internal/repository"
	"github.com/phamdinhha/event-booking-service/internal/service"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type BookingControllerInterface interface {
	CreateBooking(c *gin.Context)
	GetBooking(c *gin.Context)
	DeleteBooking(c *gin.Context)
}

type HealthCheckInterface interface {
	GetHealthCheck(c *gin.Context)
}

// Controller for testing data, should be in another service
type EventControllerInterface interface {
	CreateEvent(c *gin.Context)
	GetEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
}

type ControllerFactory struct {
	db     *sqlx.DB
	logger logger.Logger
	redis  *redis.Client
}

func NewControllerFactory(
	db *sqlx.DB,
	logger logger.Logger,
	redis *redis.Client,
) *ControllerFactory {
	return &ControllerFactory{
		db:     db,
		logger: logger,
		redis:  redis,
	}
}

func (f *ControllerFactory) NewBookingController() BookingControllerInterface {
	bookingRepo := repository.NewBookingRepository(f.db, f.logger)
	bookingSrv := service.NewBookingService(bookingRepo, f.logger, f.redis)
	return NewBookingController(f.logger, bookingSrv)
}

func (f *ControllerFactory) NewHealthCheckController() HealthCheckInterface {
	return NewHealthCheckController(f.logger, f.redis, f.db)
}

func (f *ControllerFactory) NewEventController() EventControllerInterface {
	eventRepo := repository.NewEventRepository(f.db, f.logger)
	eventSrv := service.NewEventService(eventRepo, f.logger, f.redis)
	return NewEventController(f.logger, eventSrv)
}
