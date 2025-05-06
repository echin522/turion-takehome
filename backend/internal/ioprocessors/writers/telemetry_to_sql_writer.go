package writers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"turion-takehome/internal/turiondatapacket"

	"go.uber.org/zap"
)

type TelemetryToSQLWriter struct {
	logger *zap.Logger
	db     *sql.DB // TODO: Create an interface with mocks to handle this
}

func NewTelemetryToSQLWriter(
	logger *zap.Logger,
	db *sql.DB,
) (*TelemetryToSQLWriter, error) {
	var errs error
	if logger == nil {
		errs = errors.Join(errs, errors.New("logger cannot be nil"))
	}

	if db == nil {
		errs = errors.Join(errs, errors.New("db cannot be nil"))
	}

	if errs != nil {
		return nil, errs
	}

	return &TelemetryToSQLWriter{
		logger: logger,
		db:     db,
	}, nil
}

func (w *TelemetryToSQLWriter) Write(ctx context.Context, b []byte) (int, error) {
	tdp := &turiondatapacket.TurionDataPacket{}
	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.BigEndian, tdp); err != nil {
		return 0, err
	}
	w.logger.Debug("Parsed new telemetry message", zap.Any("Message contents", tdp))

	err := insertDataPacket(ctx, w.db, tdp)
	if err != nil {
		return 0, err
	}

	w.logger.Debug("Successfully inserted TDP to DB", zap.Any("TDP contents", tdp))

	return len(b), nil
}

func (w *TelemetryToSQLWriter) Close() error {
	return nil
}

// I haven't used pure SQL in a really long time... I've been cheating and using
// Hasura as my ORM, so this probably looks stupid as hell
// insertDataPacket does a single INSERT … VALUES (…) with all 9 columns.
func insertDataPacket(
	ctx context.Context,
	db *sql.DB,
	dp *turiondatapacket.TurionDataPacket,
) error {
	const stmt = `
    INSERT INTO turion_data_packets
      (packet_id, packet_seq_ctrl, packet_length,
       ts,     subsystem_id,
       temperature, battery, altitude, signal)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.ExecContext(ctx, stmt,
		dp.CCSDSPrimaryHeader.PacketID,
		dp.CCSDSPrimaryHeader.PacketSeqCtrl,
		dp.CCSDSPrimaryHeader.PacketLength,
		dp.CCSDSSecondaryHeader.Timestamp,
		dp.CCSDSSecondaryHeader.SubsystemID,
		dp.TelemetryPayload.Temperature,
		dp.TelemetryPayload.Battery,
		dp.TelemetryPayload.Altitude,
		dp.TelemetryPayload.Signal,
	)
	return err
}
