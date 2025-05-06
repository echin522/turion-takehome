package store

import (
	"context"
	"database/sql"
	"fmt"
	"turion-takehome/internal/turiondatapacket"

	"go.uber.org/zap"
)

// DataPacketStore defines all the operations you need
// around TurionDataPacket in your HTTP routes.
type DataPacketStore interface {

	// FetchByTimeRange returns all packets whose Timestamp (ts)
	// lies in [startTS, endTS], inclusive.
	FetchByTimeRange(ctx context.Context, startTS, endTS uint64) ([]*turiondatapacket.TurionDataPacket, error)

	// FetchAnomaliesByTimeRange returns all anomalies whose Timestamp (ts)
	// lies in [startTS, endTS], inclusive.
	FetchAnomaliesByTimeRange(ctx context.Context, startTS, endTS uint64) ([]*turiondatapacket.Anomaly, error)

	// Returns the single most‐recent packet (highest ts), or ErrNoRows if none.
	FetchLatest(ctx context.Context) (turiondatapacket.TurionDataPacket, error)

	// Insert writes a new packet into the DB.
	Insert(ctx context.Context, pkt *turiondatapacket.TurionDataPacket) error

	// FetchPayloadStatsByTimeRange computes min, max, avg of each TelemetryPayload
	// column for packets whose ts is between startTS and endTS.
	FetchPayloadStatsByTimeRange(ctx context.Context, startTS, endTS uint64) (turiondatapacket.PayloadStats, error)
}

// sqlDataPacketStore is a Postgres implementation of DataPacketStore.
type sqlDataPacketStore struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewSQLDataPacketStore constructs a DataPacketStore backed by *sql.DB.
func NewSQLDataPacketStore(db *sql.DB, logger *zap.Logger) DataPacketStore {
	return &sqlDataPacketStore{db: db, logger: logger}
}

func (s *sqlDataPacketStore) FetchLatest(ctx context.Context) (turiondatapacket.TurionDataPacket, error) {
	const q = `
    SELECT packet_id, packet_seq_ctrl, packet_length,
           ts, subsystem_id,
           temperature, battery, altitude, signal
    FROM turion_data_packets
    ORDER BY ts DESC
    LIMIT 1`
	row := s.db.QueryRowContext(ctx, q)
	var ph turiondatapacket.CCSDSPrimaryHeader
	var sh turiondatapacket.CCSDSSecondaryHeader
	var tp turiondatapacket.TelemetryPayload
	if err := row.Scan(
		&ph.PacketID, &ph.PacketSeqCtrl, &ph.PacketLength,
		&sh.Timestamp, &sh.SubsystemID,
		&tp.Temperature, &tp.Battery, &tp.Altitude, &tp.Signal,
	); err != nil {
		return turiondatapacket.TurionDataPacket{}, err
	}
	return turiondatapacket.TurionDataPacket{CCSDSPrimaryHeader: ph, CCSDSSecondaryHeader: sh, TelemetryPayload: tp}, nil
}

func (s *sqlDataPacketStore) FetchByTimeRange(
	ctx context.Context,
	startTS, endTS uint64,
) ([]*turiondatapacket.TurionDataPacket, error) {
	const query = `
      SELECT packet_id, packet_seq_ctrl, packet_length,
             ts, subsystem_id,
             temperature, battery, altitude, signal
      FROM turion_data_packets
      WHERE ts BETWEEN $1 AND $2
      ORDER BY ts ASC`
	rows, err := s.db.QueryContext(ctx, query, startTS, endTS)
	if err != nil {
		return nil, fmt.Errorf("query packets by time range: %w", err)
	}
	defer rows.Close()

	// build a slice of pointers
	var packets []*turiondatapacket.TurionDataPacket
	for rows.Next() {
		var ph turiondatapacket.CCSDSPrimaryHeader
		var sh turiondatapacket.CCSDSSecondaryHeader
		var tp turiondatapacket.TelemetryPayload

		if err := rows.Scan(
			&ph.PacketID,
			&ph.PacketSeqCtrl,
			&ph.PacketLength,
			&sh.Timestamp,
			&sh.SubsystemID,
			&tp.Temperature,
			&tp.Battery,
			&tp.Altitude,
			&tp.Signal,
		); err != nil {
			return nil, fmt.Errorf("scan packet row: %w", err)
		}

		// take the address of a newly constructed TurionDataPacket
		pkt := &turiondatapacket.TurionDataPacket{
			CCSDSPrimaryHeader:   ph,
			CCSDSSecondaryHeader: sh,
			TelemetryPayload:     tp,
		}
		packets = append(packets, pkt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating packet rows: %w", err)
	}
	return packets, nil
}

func (s *sqlDataPacketStore) FetchAnomaliesByTimeRange(
	ctx context.Context,
	startTS, endTS uint64,
) ([]*turiondatapacket.Anomaly, error) {
	const q = `
      SELECT field, value, timestamp
        FROM public.anomalies
       WHERE timestamp BETWEEN $1 AND $2
       ORDER BY timestamp ASC`
	rows, err := s.db.QueryContext(ctx, q, startTS, endTS)
	if err != nil {
		return nil, fmt.Errorf("query anomalies: %w", err)
	}
	defer rows.Close()

	var out []*turiondatapacket.Anomaly
	for rows.Next() {
		var fieldName string
		var value float32
		var ts uint64

		if err := rows.Scan(&fieldName, &value, &ts); err != nil {
			return nil, fmt.Errorf("scan anomaly row: %w", err)
		}

		enumVal, found := turiondatapacket.Field_value[fieldName]
		f := turiondatapacket.Field_FIELD_UNSPECIFIED
		if found {
			f = turiondatapacket.Field(enumVal)
		}

		a := &turiondatapacket.Anomaly{
			Field:     f,
			Value:     value,
			Timestamp: ts,
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate anomaly rows: %w", err)
	}
	return out, nil
}

func (s *sqlDataPacketStore) Insert(ctx context.Context, pkt *turiondatapacket.TurionDataPacket) error {
	const stmt = `
      INSERT INTO turion_data_packets
        (packet_id, packet_seq_ctrl, packet_length,
         ts, subsystem_id,
         temperature, battery, altitude, signal)
      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := s.db.ExecContext(ctx, stmt,
		pkt.CCSDSPrimaryHeader.PacketID,
		pkt.CCSDSPrimaryHeader.PacketSeqCtrl,
		pkt.CCSDSPrimaryHeader.PacketLength,
		pkt.CCSDSSecondaryHeader.Timestamp,
		pkt.CCSDSSecondaryHeader.SubsystemID,
		pkt.TelemetryPayload.Temperature,
		pkt.TelemetryPayload.Battery,
		pkt.TelemetryPayload.Altitude,
		pkt.TelemetryPayload.Signal,
	)
	if err != nil {
		return fmt.Errorf("insert packet: %w", err)
	}
	return nil
}

func (s *sqlDataPacketStore) FetchPayloadStatsByTimeRange(
	ctx context.Context, startTS, endTS uint64,
) (turiondatapacket.PayloadStats, error) {
	const q = `
    SELECT 
      MIN(temperature), MAX(temperature), AVG(temperature),
      MIN(battery),     MAX(battery),     AVG(battery),
      MIN(altitude),    MAX(altitude),    AVG(altitude),
      MIN(signal),      MAX(signal),      AVG(signal)
    FROM turion_data_packets
    WHERE ts BETWEEN $1 AND $2`
	row := s.db.QueryRowContext(ctx, q, startTS, endTS)

	var ps turiondatapacket.PayloadStats
	// AVG returns a float64, so scan into float64 then cast
	var avgT, avgB, avgA, avgS float64

	if err := row.Scan(
		&ps.MinTemperature, &ps.MaxTemperature, &avgT,
		&ps.MinBattery, &ps.MaxBattery, &avgB,
		&ps.MinAltitude, &ps.MaxAltitude, &avgA,
		&ps.MinSignal, &ps.MaxSignal, &avgS,
	); err != nil {
		if err == sql.ErrNoRows {
			// no packets in range: return zeroes
			return turiondatapacket.PayloadStats{}, nil
		}
		return turiondatapacket.PayloadStats{}, fmt.Errorf(
			"scan payload stats: %w",
			err,
		)
	}

	// cast the averages from float64 → float32
	ps.AvgTemperature = float32(avgT)
	ps.AvgBattery = float32(avgB)
	ps.AvgAltitude = float32(avgA)
	ps.AvgSignal = float32(avgS)

	return ps, nil
}
