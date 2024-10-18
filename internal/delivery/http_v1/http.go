package http_v1

import "github.com/gin-gonic/gin"

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
