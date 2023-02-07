package server

import (
	"github.com/twoshark/alluvial1-1/common"
	"github.com/twoshark/alluvial1-1/upstream"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BalanceProxy contains the handlers for the server.
type BalanceProxy struct {
	UpstreamManager *upstream.Manager
}

func NewBalanceProxy(config common.AppConfiguration) *BalanceProxy {
	return &BalanceProxy{
		UpstreamManager: upstream.NewManager(config.Endpoints),
	}
}

func (bp BalanceProxy) RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "server root")
}
