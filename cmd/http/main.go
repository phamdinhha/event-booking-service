package main

import (
	"context"
	"log"

	"github.com/phamdinhha/event-booking-service/config"
	"github.com/phamdinhha/event-booking-service/internal/server"
	"github.com/phamdinhha/event-booking-service/pkg/db/postgres"
	"github.com/phamdinhha/event-booking-service/pkg/db/redis_client"
	"github.com/phamdinhha/event-booking-service/pkg/logger"
)

func main() {
	log.Println("Starting server...")
	cfg, err := config.GetEnvConfig()
	if err != nil {
		log.Fatalf("Error getting env config: %v", err)
	}

	var serverMode string
	if cfg.Server.Development {
		serverMode = "development"
	} else {
		serverMode = "production"
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"Server mode: %s, Port: %s, DB: %s, Redis: %s",
		serverMode,
		cfg.Server.Port,
		cfg.Postgres.Host,
		cfg.Redis.Host,
	)

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		appLogger.Fatalf("Error connecting to db: %v", err)
	}
	if err := db.Ping(); err != nil {
		appLogger.Fatalf("Error connecting to db: %v", err)
	}

	// Run migrations
	if err := postgres.RunMigrations(db, cfg.Migrations.Path); err != nil {
		appLogger.Fatalf("Error running migrations: %v", err)
	}

	redisClient := redis_client.NewRedisClient(cfg)
	if err := redisClient.Ping(context.TODO()).Err(); err != nil {
		appLogger.Fatalf("Error connecting to redis: %v", err)
	}
	// Apply migrations
	//
	server := server.NewServer(appLogger, cfg, redisClient, db)
	appLogger.Info("Starting server...")
	appLogger.Fatal(server.Run(context.Background()))
	appLogger.Info("Server started")
}
