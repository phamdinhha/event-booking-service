package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	srvAddr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)

	handlers := s.SetupHandlers()
	server := &http.Server{
		Addr:    srvAddr,
		Handler: handlers,
	}

	// Create a context that will be canceled on interrupt signal
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// Start the server in a goroutine
	go func() {
		s.logger.Info("Server is running on " + srvAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	s.logger.Info("Server is shutting down")

	// Create a timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdown = func() {
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Server shutdown failed", "error", err)
		}
	}

	return shutdown, nil
}

func (s *Server) SetupHandlers() *gin.Engine {
	ginEngine := gin.Default()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(gin.Logger())

	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	return ginEngine
}

func (s *Server) MapHandlers(ginEngine *gin.Engine) {

}
