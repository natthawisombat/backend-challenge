package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func StartUserCountLogger(ctx context.Context, db *mongo.Database, logger *zap.SugaredLogger) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		coll := db.Collection("user")

		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopped user count logger")
				return
			case <-ticker.C:
				count, err := coll.CountDocuments(ctx, bson.M{})
				if err != nil {
					logger.Errorw("Failed to count users", "error", err)
					continue
				}
				logger.Infow("Current user count", "count", count)
			}
		}
	}()
}
