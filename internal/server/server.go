package server

import (
	"github.com/phamdinhha/event-booking-service/config"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	logger logger.Logger
	cfg    *config.Config
	redis  *redis.Client
}

func NewServer(
	logger logger.Logger,
	cfg *config.Config,
	redis *redis.Client,
)
