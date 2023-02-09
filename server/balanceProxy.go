package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// BalanceProxy contains the handlers for the server.
type BalanceProxy struct{}

func NewBalanceProxy() *BalanceProxy {
	return new(BalanceProxy)
}

func (bp BalanceProxy) RootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "server root")
}
