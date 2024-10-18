package http_v1

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/phamdinhha/event-booking-service/internal/dto"
	"github.com/phamdinhha/event-booking-service/internal/service"
	"github.com/phamdinhha/event-booking-service/pkg/http_utils"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
)

type EventController struct {
	logger   logger.Logger
	eventSrv service.EventServiceInterface
}

func NewEventController(
	logger logger.Logger,
	eventSrv service.EventServiceInterface,
) EventControllerInterface {
	return &EventController{logger: logger, eventSrv: eventSrv}
}

func (e *EventController) CreateEvent(c *gin.Context) {
	var req dto.CreateEventDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		e.logger.Error("EVENT_CONTROLLER.CREATE_EVENT.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(http_utils.INVALID_REQUEST, err.Error()))
		return
	}
	created, err := e.eventSrv.CreateEvent(c.Request.Context(), &req)
	if err != nil {
		e.logger.Error("EVENT_CONTROLLER.CREATE_EVENT.Error", err)
		c.JSON(http.StatusInternalServerError, http_utils.NewErrorResponse(http_utils.INTERNAL_SERVER_ERROR, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, http_utils.NewOKResponse(http_utils.CREATED, created))
}

func (e *EventController) GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		e.logger.Error("EVENT_CONTROLLER.GET_EVENT.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(
			http_utils.INVALID_REQUEST,
			err,
		))
		return
	}

	event, err := e.eventSrv.GetEventByID(c.Request.Context(), id)
	if err != nil {
		e.logger.Error("EVENT_CONTROLLER.GET_EVENT.Error", err)
		var statusCode int
		var message string
		if err == sql.ErrNoRows {
			statusCode = http.StatusNotFound
			message = http_utils.NOT_FOUND
		} else {
			statusCode = http.StatusInternalServerError
			message = http_utils.INTERNAL_SERVER_ERROR
		}
		c.JSON(statusCode, http_utils.NewErrorResponse(message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, http_utils.NewOKResponse(
		http_utils.SUCCESS,
		event,
	))
}

func (e *EventController) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	id, err := uuid.Parse(eventID)
	if err != nil {
		e.logger.Error("EVENT_CONTROLLER.DELETE_EVENT.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(
			http_utils.INVALID_REQUEST,
			err,
		))
		return
	}
	err = e.eventSrv.DeleteEvent(c.Request.Context(), id)
	if err != nil {
		e.logger.Error("EVENT_CONTROLLER.DELETE_EVENT.Error", err)
		var statusCode int
		var message string
		if err == sql.ErrNoRows {
			statusCode = http.StatusNotFound
			message = http_utils.NOT_FOUND
		} else {
			statusCode = http.StatusInternalServerError
			message = http_utils.INTERNAL_SERVER_ERROR
		}
		response := http_utils.NewErrorResponse(message, err.Error())
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, http_utils.NewOKResponse(
		http_utils.SUCCESS,
		gin.H{"message": "Event successfully deleted"},
	))
}
