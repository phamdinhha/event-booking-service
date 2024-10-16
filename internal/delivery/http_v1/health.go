package http_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type HealthCheckController struct {
	logger logger.Logger
	redis  *redis.Client
	db     *sqlx.DB
}

func NewHealthCheckController(
	logger logger.Logger,
	redis *redis.Client,
	db *sqlx.DB,
) HealthCheckInterface {
	return &HealthCheckController{logger: logger, redis: redis, db: db}
}

func (h *HealthCheckController) GetHealthCheck(c *gin.Context) {
	status := "healthy"
	statusCode := http.StatusOK

	// Check Redis connection
	_, err := h.redis.Ping(c).Result()
	if err != nil {
		h.logger.Error("Redis health check failed", "error", err)
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	// Check DB connection
	err = h.db.Ping()
	if err != nil {
		h.logger.Error("Database health check failed", "error", err)
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status": status,
	})
}
