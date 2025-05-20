package main

import (
	"backend-challenge/configs"
	"backend-challenge/pkg/logging"
	"backend-challenge/routers"
	"backend-challenge/utils"
	"context"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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

	// task background process
	utils.StartUserCountLogger(ctx, app.DBMongo.DB, logger)

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
