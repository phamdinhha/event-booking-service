package http_v1

import "github.com/gin-gonic/gin"

func MapHealthCheckRoutes(
	router *gin.RouterGroup,
	controller HealthCheckInterface,
) {
	router.GET("/", controller.GetHealthCheck)
}

func MapBookingRoutes(
	router *gin.RouterGroup,
	controller BookingControllerInterface,
) {
	router.POST("/", controller.CreateBooking)
	router.GET("/:id", controller.GetBooking)
	router.DELETE("/:id", controller.DeleteBooking)
}

func MapEventRoutes(
	router *gin.RouterGroup,
	controller EventControllerInterface,
) {
	router.POST("/", controller.CreateEvent)
	router.GET("/:id", controller.GetEvent)
	router.DELETE("/:id", controller.DeleteEvent)
}
