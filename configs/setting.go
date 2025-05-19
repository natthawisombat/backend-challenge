package configs

import (
	"backend-challenge/configs/store"
	"backend-challenge/middlewares"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Setting struct {
	App     *fiber.App
	Logger  *zap.SugaredLogger
	DBMongo *store.MongoStore
}

func NewApp(logger *zap.SugaredLogger) *Setting {
	return &Setting{Logger: logger}
}

func (c *Setting) SetApp(ctx context.Context) error {
	c.Logger.Named("backend-chellenge")

	c.App = fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	if err := SetEnv(ctx); err != nil {
		return err
	}

	cfg := cors.Config{
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}

	c.App.Use(recover.New(recover.Config{EnableStackTrace: true}))
	c.App.Use(helmet.New())
	c.App.Use(cors.New(cfg))
	c.App.Use(middlewares.LoggerMiddleware(c.Logger))
	c.App.Use(logger.New(logger.Config{
		Format:     "${blue}${time} ${yellow}${status} - ${red}${latency} ${cyan}${method} ${path} ${green} ${ip} ${ua} ${reset}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Bangkok",
		Output:     os.Stdout,
	}))

	mongodb, err := store.ConnectMongo(ctx)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Mongo connect failed error %s ", err.Error())
	}

	c.DBMongo = mongodb
	return nil
}

func (c *Setting) RunApp(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		err := c.App.Listen(fmt.Sprintf(":%v", App.Port))
		if err != nil {
			errChan <- fmt.Errorf("fiber listen error: %w", err)
		}
	}()

	return errChan
}

func (c *Setting) StopApp(ctx context.Context) error {
	if err := c.App.Shutdown(); err != nil {
		return fmt.Errorf("error app failed to stop : %s", err)
	}

	// Disconnect Mongo
	if c.DBMongo != nil {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := c.DBMongo.Client.Disconnect(ctxTimeout); err != nil {
			return fmt.Errorf("mongo disconnect error: %w", err)
		}
	}

	return nil
}
