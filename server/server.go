package server

import (
	"context"
	"errors"
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

func Start(config common.AppConfiguration) {
	var wg sync.WaitGroup
	wg.Add(1)

	bp := NewBalanceProxy(config)
	go func() {
		defer wg.Done()
		bp.InitClients(config)
	}()
	wg.Wait()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", bp.RootHandler)
	e.GET("/ethereum/balance/:account", bp.GetLatestBalance)
	e.GET("/ethereum/balance/:account/block/:block", bp.GetBalance)

	// Start server
	if err := e.Start(":" + config.ListenPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("shutting down the server")
	}
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	timeout := viper.GetInt("SHUTDOWN_TIMEOUT")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
