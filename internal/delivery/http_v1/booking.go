package http_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/phamdinhha/event-booking-service/internal/service"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
)

type BookingController struct {
	logger     logger.Logger
	bookingSrv service.BookingService
}

func NewBookingController(
	logger logger.Logger,
	bookingSrv service.BookingService,
) BookingControllerInterface {
	return &BookingController{logger: logger, bookingSrv: bookingSrv}
}

func (b *BookingController) CreateBooking(c *gin.Context) {

}

func (b *BookingController) GetBooking(c *gin.Context) {

}

func (b *BookingController) DeleteBooking(c *gin.Context) {

}
