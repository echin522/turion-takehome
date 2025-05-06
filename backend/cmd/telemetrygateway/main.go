package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"turion-takehome/internal/config"
	"turion-takehome/internal/ioprocessors"
	"turion-takehome/internal/ioprocessors/quarantiners"
	"turion-takehome/internal/ioprocessors/readers"
	"turion-takehome/internal/ioprocessors/writers"
	"turion-takehome/internal/utils"

	_ "github.com/jackc/pgx/v5/stdlib" // register the "pgx" driver
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	TDP_BUFFER_SIZE int = 1024
)

// The primary entrypoint for telemetry
// 1. Opens up a UDP client to receive data from the telemetry generator
// 2. Deserializes the UDP packet and extracts the telemetry payload
// 3. Checks for any anomalies
// 4. Persists the data in a PG db
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devMode := flag.Bool("loc", false, "Run in local development mode")
	flag.Parse()

	logger := utils.InitLogger()
	if *devMode {
		logger = utils.InitDevLogger()
	}

	envConfig, err := config.NewTelemetryGatewayConfig()
	if err != nil {
		logger.Fatal("Failed to create new telemetry gateway config", zap.Error(err))
	}

	eg, ctx := errgroup.WithContext(ctx)

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

	err = applyMigrations(db, "/migrations")
	if err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Channels for processes
	anomalyChannel := make(chan []byte)
	defer close(anomalyChannel)
	sqlChannel := make(chan []byte)
	defer close(sqlChannel)

	anomalyChannelWriter := writers.NewChannelWriter(logger, anomalyChannel)
	sqlChannelWriter := writers.NewChannelWriter(logger, sqlChannel)

	groundStationEmulatorUDPAddr, err := net.ResolveUDPAddr(
		"udp",
		envConfig.GroundStationEmulatorAddress,
	)
	if err != nil {
		logger.Fatal("Failed to create mock GS UDP server address", zap.Error(err))
	}

	tdpReader, err := readers.NewUDPReader(
		logger,
		readers.WithUDPAddr(groundStationEmulatorUDPAddr),
	)
	if err != nil {
		logger.Fatal("Failed to create telemetry UDP reader", zap.Error(err))
	}

	// I would normally put quarantined messages in a different table or in some
	// other quarantine zone to attempt reprocessing/manual debugging, but due to
	// time constraints I will just log it
	tdpQuarantiner, err := quarantiners.NewNoOpQuarantiner(logger)
	if err != nil {
		logger.Fatal("Failed to create new quarantiner", zap.Error(err))
	}

	tdpWriter, err := writers.NewTelemetryMessageWriter(
		logger,
		sqlChannelWriter,
		anomalyChannelWriter,
	)
	if err != nil {
		logger.Fatal("Failed to create new Turion Data Packet writer", zap.Error(err))
	}

	telemetryProcessor := ioprocessors.NewProcessor(
		logger,
		TDP_BUFFER_SIZE,
		tdpReader, tdpWriter, tdpQuarantiner,
	)
	defer func() {
		if err := telemetryProcessor.Close(); err != nil {
			logger.Fatal("Failed to close telemetry to SQL processor", zap.Error(err))
		}
	}()
	eg.Go(func() error {
		err := telemetryProcessor.Start(ctx)
		if err != nil {
			logger.Error("Param processor encountered an error", zap.Error(err))
			return err
		}
		return nil
	})
	logger.Info("Param processor started")

	sqlChannelReader := readers.NewChannelReader(logger, sqlChannel)
	sqlWriter, err := writers.NewTelemetryToSQLWriter(logger, db)
	if err != nil {
		logger.Fatal("Failed to create new SQL writer", zap.Error(err))
	}
	sqlQuarantiner, err := quarantiners.NewNoOpQuarantiner(logger)
	if err != nil {
		logger.Fatal("Failed to create new quarantiner", zap.Error(err))
	}
	telemetryToSQLProcessor := ioprocessors.NewProcessor(
		logger,
		TDP_BUFFER_SIZE,
		sqlChannelReader, sqlWriter, sqlQuarantiner,
	)
	defer func() {
		if err := telemetryToSQLProcessor.Close(); err != nil {
			logger.Fatal("Failed to close telemetry to SQL processor", zap.Error(err))
		}
	}()
	eg.Go(func() error {
		err := telemetryToSQLProcessor.Start(ctx)
		if err != nil {
			logger.Error("Telemetry to SQL processor failed to start", zap.Error(err))
			return err
		}
		return nil
	})

	anomalyChannelReader := readers.NewChannelReader(logger, anomalyChannel)
	anomalyWriter, err := writers.NewAnomalyWriter(logger, db)
	if err != nil {
		logger.Fatal("Failed to create new anomaly writer", zap.Error(err))
	}
	anomalyQuarantiner, err := quarantiners.NewNoOpQuarantiner(logger)
	if err != nil {
		logger.Fatal("Failed to create new quarantiner", zap.Error(err))
	}
	anomalyProcessor := ioprocessors.NewProcessor(
		logger,
		TDP_BUFFER_SIZE,
		anomalyChannelReader, anomalyWriter, anomalyQuarantiner,
	)
	defer func() {
		if err := anomalyProcessor.Close(); err != nil {
			logger.Fatal("Failed to close anomaly processor", zap.Error(err))
		}
	}()
	eg.Go(func() error {
		err := anomalyProcessor.Start(ctx)
		if err != nil {
			logger.Error("Anomaly processor failed to start", zap.Error(err))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logger.Error("Some process has terminated", zap.Error(err))
	}
	logger.Info("All services have stopped")
}

// I couldn't get the golang docker migrator to work so I threw this together instead
func applyMigrations(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}
	// ensure lex order
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" && strings.HasSuffix(e.Name(), ".up.sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join(dir, name)
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		if _, err := db.Exec(string(b)); err != nil {
			return fmt.Errorf("exec %s: %w", name, err)
		}
		fmt.Printf("applied migration %s\n", name)
	}
	return nil
}
