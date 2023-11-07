package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	r.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "/ping")
	})
	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world")
	})
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// Start server
	go func() {
		if err := r.Start(":8080"); err != nil {
			r.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}
}
