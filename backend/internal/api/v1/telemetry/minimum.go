package telemetry

import (
	"net/http"
	"time"
	"turion-takehome/internal/store"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func AggregationHandler(store store.DataPacketStore, logger *zap.Logger) func(echo.Context) error {
	return func(c echo.Context) error {
		startStr := c.QueryParam("start_time")
		endStr := c.QueryParam("end_time")
		if startStr == "" || endStr == "" {
			return echo.NewHTTPError(http.StatusBadRequest,
				"`start_time` and `end_time` are required (ISO8601)")
		}

		startT, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest,
				"invalid start_time: "+err.Error())
		}
		endT, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest,
				"invalid end_time: "+err.Error())
		}

		startTS := uint64(startT.Unix())
		endTS := uint64(endT.Unix())

		stats, err := store.FetchPayloadStatsByTimeRange(c.Request().Context(), startTS, endTS)
		if err != nil {
			logger.Error("failed to fetch payload stats", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, stats)
	}
}
