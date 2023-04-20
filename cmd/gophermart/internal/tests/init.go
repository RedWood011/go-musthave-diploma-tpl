package tests

import (
	"context"

	"github.com/RedWood011/cmd/gophermart/internal/config"
	"github.com/RedWood011/cmd/gophermart/internal/database/postgres"
	"github.com/RedWood011/cmd/gophermart/internal/logger"
	"github.com/RedWood011/cmd/gophermart/internal/service"
	"github.com/RedWood011/cmd/gophermart/internal/transport/http"
	"github.com/gofiber/fiber/v2"
)

func initTest() (*fiber.App, error) {
	ctx := context.Background()
	log := logger.InitLogger()
	cfg := config.New()
	db, err := postgres.NewDatabase(ctx, cfg.DataBaseURI, cfg.CountRepetitionBD)
	if err != nil {

		return nil, err
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Info("Failed to ping to database")
		return nil, err
	}

	serviceHTTP := service.NewService(db, cfg, log)

	serverParam := http.ServerParams{
		Service: serviceHTTP,
		Storage: db,
		Cfg:     cfg,
		Logger:  log,
	}
	server := http.NewServer(serverParam)

	return server, nil
}
