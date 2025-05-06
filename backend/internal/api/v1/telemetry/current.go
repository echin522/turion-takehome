package telemetry

import (
	"database/sql"
	"net/http"
	"turion-takehome/internal/store"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func CurrentHandler(store store.DataPacketStore, logger *zap.Logger) func(echo.Context) error {
	return func(c echo.Context) error {
		pkt, err := store.FetchLatest(c.Request().Context())
		if err != nil {
			if err == sql.ErrNoRows {
				return c.NoContent(http.StatusNoContent)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, pkt)
	}
}
