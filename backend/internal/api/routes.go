package api

import (
	"net/http"
	telemhandlers "turion-takehome/internal/api/v1/telemetry"
	"turion-takehome/internal/store"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	ROUTE_PING                   = "/api/v1/ping"
	ROUTE_TELEMETRY              = "/api/v1/telemetry"
	ROUTE_TELEMETRY_CURRENT      = "/api/v1/telemetry/current"
	ROUTE_TELEMETRY_ANOMALIES    = "/api/v1/telemetry/anomaly"
	ROUTE_TELEMETRY_AGGREGATIONS = "/api/v1/telemetry/aggregation"
	ROUTE_ANOMALIES_NEW          = "/api/v1/anomaly/new"
)

// RegisterRoutes mounts all of your telemetry routes onto the Echo instance.
func RegisterRoutes(e *echo.Echo, store store.DataPacketStore, logger *zap.Logger) {
	// Basic health check
	e.GET(ROUTE_PING, func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// GET /api/v1/telemetry?start_time=<ISO>&end_time=<ISO>
	e.GET(ROUTE_TELEMETRY, telemhandlers.TelemetryList(store, logger))

	// GET /api/v1/telemetry/current
	e.GET(ROUTE_TELEMETRY_CURRENT, telemhandlers.CurrentHandler(store, logger))

	// GET /api/v1/telemetry/anomalies?start_time=<ISO>&end_time=<ISO>
	e.GET(ROUTE_TELEMETRY_ANOMALIES, telemhandlers.AnomaliesHandler(store, logger))

	// GET /api/v1/telemetry/stats?start_time=<ISO>&end_time=<ISO>
	e.GET(ROUTE_TELEMETRY_AGGREGATIONS, telemhandlers.AggregationHandler(store, logger))
}
