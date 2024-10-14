package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/phamdinhha/event-booking-service/config"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/phamdinhha/event-booking-service/pkg/utils"
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
) *Server {
	return &Server{
		logger: logger,
		cfg:    cfg,
		redis:  redis,
	}
}

func (s *Server) Run(ctx context.Context) (shutdown utils.Deamon, err error) {
	var (
		srvAddr = fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	)
	// handlers := handlers.NewHandlers(s.logger, s.redis)

	server := &http.Server{
		Addr: srvAddr,
		// Handler: handlers
	}
	s.logger.Info("Server is running on", srvAddr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server is not running", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.logger.Info("Server is shutting down")

	shutdown = func() {
		<-ctx.Done()
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Fatal("Server is not shutting down", err)
		}
	}

	return shutdown, nil
}
