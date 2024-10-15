package http

import "github.com/gin-gonic/gin"

// type BookingController interface {

// }

type HealthCheckInterface interface {
	GetHealthCheck(c *gin.Context)
}
