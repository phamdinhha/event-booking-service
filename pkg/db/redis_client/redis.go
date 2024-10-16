package redis_client

import (
	"fmt"

	"github.com/phamdinhha/event-booking-service/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(c *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
}
