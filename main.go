package main

import (
	"backend-challenge/configs"
	"backend-challenge/pkg/logging"
	"backend-challenge/routers"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)
	defer stop()

	// app
	app := configs.NewApp(logger)
	if err := app.SetApp(ctx); err != nil {
		logger.Fatal(err)
	}

	routers.SetupRoutes(app)
	errChan := app.RunApp(ctx)

	select {
	case <-ctx.Done():
		logger.Info("shutting down via signal...")
	case err := <-errChan:
		logger.Errorw("server error", "error", err)
	}

	if err := app.StopApp(ctx); err != nil {
		logger.Errorw("shutdown error", "error", err)
	} else {
		logger.Info("shutdown complete")
	}
}
