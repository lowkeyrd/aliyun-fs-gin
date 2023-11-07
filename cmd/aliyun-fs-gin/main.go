package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// Defining the Graphql handler
/*
func graphqlHandler(logger *zap.Logger) echo.HandlerFunc {
	srv := handler.NewDefaultServer()
	return func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
*/

func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/")
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := echo.New()

	r.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	r.GET("/", playgroundHandler())
	//r.POST("/", graphqlHandler(client, logger))
	//r.POST("/", graphqlHandler(logger))
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
