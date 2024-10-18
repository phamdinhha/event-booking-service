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

type BookingController struct {
	logger     logger.Logger
	bookingSrv service.BookingServiceInterface
}

func NewBookingController(
	logger logger.Logger,
	bookingSrv service.BookingServiceInterface,
) BookingControllerInterface {
	return &BookingController{logger: logger, bookingSrv: bookingSrv}
}

func (b *BookingController) CreateBooking(c *gin.Context) {
	var req dto.CreateBookingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		b.logger.Error("BOOKING_CONTROLLER.CREATE_BOOKING.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(http_utils.INVALID_REQUEST, err.Error()))
		return
	}
	created, err := b.bookingSrv.CreateBooking(c.Request.Context(), &req)
	if err != nil {
		b.logger.Error("BOOKING_CONTROLLER.CREATE_BOOKING.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(http_utils.INTERNAL_SERVER_ERROR, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, http_utils.NewOKResponse(http_utils.CREATED, created))
}

func (b *BookingController) GetBooking(c *gin.Context) {
	bookingID := c.Param("id")
	id, err := uuid.Parse(bookingID)
	if err != nil {
		b.logger.Error("BOOKING_CONTROLLER.GET_BOOKING.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(
			http_utils.INVALID_REQUEST,
			err,
		))
		return
	}

	booking, err := b.bookingSrv.GetBooking(c.Request.Context(), id)
	if err != nil {
		b.logger.Error("BOOKING_CONTROLLER.GET_BOOKING.Error", err)
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
		booking,
	))
}

func (b *BookingController) DeleteBooking(c *gin.Context) {
	bookingID := c.Param("id")
	id, err := uuid.Parse(bookingID)
	if err != nil {
		b.logger.Error("BOOKING_CONTROLLER.GET_BOOKING.Error", err)
		c.JSON(http.StatusBadRequest, http_utils.NewErrorResponse(
			http_utils.INVALID_REQUEST,
			err,
		))
		return
	}
	err = b.bookingSrv.DeleteBooking(c.Request.Context(), id)
	if err != nil {
		b.logger.Error("BOOKING_CONTROLLER.DELETE_BOOKING.Error", err)
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
		gin.H{"message": "Booking successfully deleted"},
	))

}
