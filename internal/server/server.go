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
	"github.com/jmoiron/sqlx"
	"github.com/phamdinhha/event-booking-service/config"
	"github.com/phamdinhha/event-booking-service/internal/delivery/http_v1"
	"github.com/phamdinhha/event-booking-service/internal/repository"
	"github.com/phamdinhha/event-booking-service/internal/service"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
	"github.com/phamdinhha/event-booking-service/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	logger logger.Logger
	cfg    *config.Config
	redis  *redis.Client
	db     *sqlx.DB
}

func NewServer(
	logger logger.Logger,
	cfg *config.Config,
	redis *redis.Client,
	db *sqlx.DB,
) *Server {
	return &Server{
		logger: logger,
		cfg:    cfg,
		redis:  redis,
		db:     db,
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
	s.MapHandlers(ginEngine)
	return ginEngine
}

func (s *Server) MapHandlers(ginEngine *gin.Engine) {
	healthCheckController := http_v1.NewHealthCheckController(s.logger, s.redis, s.db)
	healthCheckGroup := ginEngine.Group("/health")
	http_v1.MapHealthCheckRoutes(healthCheckGroup, healthCheckController)

	bookingRepo := repository.NewBookingRepository(s.db, s.logger)
	bookingSrv := service.NewBookingService(bookingRepo, s.logger, s.redis)
	bookingController := http_v1.NewBookingController(s.logger, bookingSrv)
	bookingGroup := ginEngine.Group("/bookings")
	http_v1.MapBookingRoutes(bookingGroup, bookingController)

	eventRepo := repository.NewEventRepository(s.db, s.logger)
	eventSrv := service.NewEventService(eventRepo, s.logger, s.redis)
	eventController := http_v1.NewEventController(s.logger, eventSrv)
	eventGroup := ginEngine.Group("/events")
	http_v1.MapEventRoutes(eventGroup, eventController)
}
