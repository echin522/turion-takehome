package writers

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"turion-takehome/internal/turiondatapacket"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type TelemetryMessageWriter struct {
	logger        *zap.Logger
	sqlWriter     Writer
	anomalyWriter Writer
}

func NewTelemetryMessageWriter(
	logger *zap.Logger,
	sqlWriter Writer,
	anomalyWriter Writer,
) (*TelemetryMessageWriter, error) {
	var errs error
	if logger == nil {
		errs = errors.Join(errors.New("logger cannot be nil"))
	}

	if sqlWriter == nil {
		errs = errors.Join(errors.New("sql writer cannot be nil"))
	}

	if anomalyWriter == nil {
		errs = errors.Join(errors.New("anomaly writer cannot be nil"))
	}

	if errs != nil {
		return nil, errs
	}

	return &TelemetryMessageWriter{
		logger:        logger,
		anomalyWriter: anomalyWriter,
		sqlWriter:     sqlWriter,
	}, nil
}

func (w *TelemetryMessageWriter) Write(ctx context.Context, b []byte) (int, error) {
	tdp := &turiondatapacket.TurionDataPacket{}
	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.BigEndian, tdp); err != nil {
		return 0, err
	}
	w.logger.Debug("Parsed new telemetry message", zap.Any("Message contents", tdp))

	writtenByteCount, err := w.sqlWriter.Write(ctx, b)
	if err != nil {
		return 0, err
	}

	anomalies := tdp.DetectAnomalies()

	// Can't use _, anomaly := range anomalies because Anomaly is a protoc message
	// Can't copy mutex
	for i := range anomalies {
		anomaly := &anomalies[i]
		a, err := proto.Marshal(anomaly)
		if err != nil {
			w.logger.Error(
				"Failed to marshal anomaly to JSON",
				zap.String("Anomaly", anomaly.Field.String()),
				zap.Error(err),
			)
		}
		awb, err := w.anomalyWriter.Write(ctx, a)
		if err != nil {
			return 0, err
		}
		writtenByteCount += awb
	}

	return (writtenByteCount), nil
}

func (w *TelemetryMessageWriter) Close() error {
	return nil
}
