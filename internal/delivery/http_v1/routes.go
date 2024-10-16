package http_v1

import "github.com/gin-gonic/gin"

func MapHealthCheckRoutes(
	router *gin.RouterGroup,
	controller HealthCheckInterface,
) {
	router.GET("/", controller.GetHealthCheck)
}
