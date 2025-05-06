package telemetry

import (
	"net/http"
	"time"
	"turion-takehome/internal/store"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func TelemetryList(
	store store.DataPacketStore,
	logger *zap.Logger,
) func(echo.Context) error {
	return func(c echo.Context) error {
		// 1) Parse query params
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

		// 2) Convert to uint64 seconds
		startTS := uint64(startT.Unix())
		endTS := uint64(endT.Unix())

		// 3) Fetch from store
		packets, err := store.FetchByTimeRange(c.Request().Context(), startTS, endTS)
		if err != nil {
			// if it's a SQL error you might inspect it, but 500 is fine
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// 4) Return JSON array of TurionDataPacket
		return c.JSON(http.StatusOK, packets)
	}
}
