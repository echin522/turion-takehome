package writers

import (
	"context"
	"database/sql"
	"errors"
	"turion-takehome/internal/turiondatapacket"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// AnomalyWriter shall inspect telemetry messages and check for any values that
// meet or exceed "safety thresholds". These thresholds are defined in the TDP
// package /internal/turiondatapacket/anomaly.go
//
// If a threshold is exceeded, an alert is created to notify users of the
// occurance. Due to time constraints, this will come in the form of an HTTP
// POST request. The route will simply contain a loud logger for the alert.
// Were this a production application, this alert would go to something like
// PagerDuty or Grafana instead.
type AnomalyWriter struct {
	logger *zap.Logger
	db     *sql.DB // TODO: Create an interface with mocks to handle this

}

func NewAnomalyWriter(
	logger *zap.Logger,
	db *sql.DB,

) (*AnomalyWriter, error) {
	var errs error
	if logger == nil {
		return nil, errors.Join(errs, errors.New("logger cannot be nil"))
	}

	if db == nil {
		errs = errors.Join(errs, errors.New("db cannot be nil"))
	}
	if errs != nil {
		return nil, errs
	}

	return &AnomalyWriter{
		logger: logger,
		db:     db,
	}, nil
}

// I wanted to use protobufs but I am now stuck because I'm trying to POST protobufs
// instead of JSON. So let's just convert it to JSON for now. Hate this, it's
// due to my lack of foresight and time constraints
//
// I should have done gRPC but I don't know how to do it with ECHO
// Or I should've just set up opentelemetry from the getgo
func (w *AnomalyWriter) Write(ctx context.Context, b []byte) (int, error) {

	var a turiondatapacket.Anomaly
	if err := proto.Unmarshal(b, &a); err != nil {
		return 0, err
	}

	// 2) Insert into SQL
	const stmt = `
    INSERT INTO public.anomalies (field, value, timestamp)
    VALUES ($1, $2, $3)`
	_, err := w.db.ExecContext(ctx, stmt,
		a.Field.String(),
		a.Value,
		a.Timestamp,
	)
	if err != nil {
		w.logger.Error("failed to insert anomaly", zap.Error(err))
		return 0, err
	}

	w.logger.Debug("wrote anomaly to database",
		zap.String("field", a.Field.String()),
		zap.Float32("value", a.Value),
		zap.Uint64("ts", a.Timestamp),
	)

	return len(b), nil
}

// TODO
func (w *AnomalyWriter) Close() error {
	return nil
}
