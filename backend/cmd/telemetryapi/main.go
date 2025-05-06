package main

import (
	"context"
	"database/sql"
	"turion-takehome/internal/api"
	"turion-takehome/internal/config"
	"turion-takehome/internal/store"
	"turion-takehome/internal/utils"

	_ "github.com/jackc/pgx/v5/stdlib" // register the "pgx" driver
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := utils.InitLogger()

	envConfig, err := config.NewTelemetryAPIConfig()
	if err != nil {
		logger.Fatal("Failed to create new telemetry API server config", zap.Error(err))
	}

	e := echo.New()

	db, err := sql.Open("pgx", envConfig.PGHostURL)
	if err != nil {
		logger.Fatal(
			"Failed to connect to PG database",
			zap.String("Host URL", envConfig.PGHostURL),
			zap.Error(err),
		)
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("✅ connected to PG db", zap.String("DB URL", envConfig.PGHostURL))

	store := store.NewSQLDataPacketStore(db, logger)

	api.RegisterRoutes(e, store, logger)

	e.Logger.Fatal(e.Start(":8090"))
}
