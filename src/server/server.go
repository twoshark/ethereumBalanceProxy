package server

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4/middleware"
	"github.com/twoshark/balanceproxy/src/common"
	"github.com/twoshark/balanceproxy/src/metrics"
	"net/http"
	"os"
	"os/signal"
	"time"

	echoProm "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Start(config common.AppConfiguration, ready chan bool) {
	startTime := time.Now()

	bp := NewBalanceProxy(config)
	bp.InitClients()

	quitHealthChecker := bp.UpstreamManager.StartHealthCheck()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", bp.RootHandler)
	e.GET("/ethereum/balance/:account", bp.GetLatestBalance)
	e.GET("/ethereum/balance/:account/block/:block", bp.GetBalance)
	e.GET("/live", bp.LiveHandler)
	e.GET("/ready", bp.ReadyHandler)

	p := echoProm.NewPrometheus("echo", nil)
	p.Use(e)

	go func() { ready <- true }()
	startupTime := time.Since(startTime)
	metrics.Metrics().StartUpTime.With(nil).Observe(float64(startupTime.Milliseconds()))
	// Start server
	if err := e.Start(":" + config.ListenPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("shutting down the server")
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	go func() { quitHealthChecker <- true }()
	timeout := viper.GetInt("SHUTDOWN_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
