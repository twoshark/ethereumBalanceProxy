package server

import (
	"context"
	"errors"
	echoProm "github.com/labstack/echo-contrib/prometheus"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/twoshark/balanceproxy/common"
)

func Start(config common.AppConfiguration, ready chan bool) {
	var wg sync.WaitGroup
	wg.Add(1)

	bp := NewBalanceProxy(config)
	bp.InitClients()

	quitHealthChecker := bp.UpstreamManager.StartHealthCheck()

	// Echo instance
	e := echo.New()

	p := echoProm.NewPrometheus("echo", nil)
	p.Use(e)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", bp.RootHandler)
	e.GET("/ethereum/balance/:account", bp.GetLatestBalance)
	e.GET("/ethereum/balance/:account/block/:block", bp.GetBalance)

	go func() { ready <- true }()

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
