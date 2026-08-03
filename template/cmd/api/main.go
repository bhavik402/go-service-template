// Command api is the entry point. It stays intentionally small: build the
// application graph in dependencies.go, then run it. main.go should never
// grow business logic.
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app, err := BuildApplication()
	if err != nil {
		log.Fatalf("failed to build application: %v", err)
	}
	defer app.Close()

	srv := &http.Server{
		Addr:              ":" + app.Config.HTTPPort,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		app.Logger.Info("starting server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	app.Logger.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Error("graceful shutdown failed", slog.Any("error", err))
	}
}
