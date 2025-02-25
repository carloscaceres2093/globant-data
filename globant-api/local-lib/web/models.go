package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"globant-api/local-lib/log"
)

const (
	port8080 = ":8080"
)

// Error represents an error in this api.
type Error struct {
	Message string `json:"message"`
}

// App contains everything you need to create an api.
type App struct {
	*Router
	Environment

	Log log.Logger

	server *http.Server
}

// Environment represents the environment of api.
type Environment struct {
	Name string
}

// Run execute listen and serve of http.
func (wa *App) Run() error {
	wa.Log.Info(context.Background(), "", log.Field("port", port8080),
		log.Field("env", wa.Name))
	wa.server = &http.Server{Addr: port8080, Handler: wa.mux}

	shutdownComplete := handleShutdown(wa)

	err := wa.server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			wa.Log.Info(context.Background(), "", log.Field("error", err.Error()))
			return err
		}
		// in case server is closed or shutdown we will not return an error because we take it as a successful end for the application
		<-shutdownComplete
	}

	return nil
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections.
// Shutdown takes as argument a context with timeout. The timeout is the maximum allowed for the current requests to complete.
func (wa *App) Shutdown(ctx context.Context) error {
	wa.Log.Info(context.Background(), "", log.Field("message", port8080))

	return wa.server.Shutdown(ctx)
}

func handleShutdown(r *App) <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

		<-shutdownSignal
		r.Log.Debug(context.Background(), "",
			log.Field("message", "received signal to shutdown application"))

		// pods some time to stop delivering requests to server
		time.Sleep(60 * time.Second)

		toCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer func(logger log.Logger) {
			cancel()
		}(r.Log)

		if err := r.Shutdown(toCtx); err != nil {
			r.Log.Warning(context.Background(), "", log.Field("message",
				fmt.Sprintf("server.Shutdown failed: %v\n", err)))
		}

		close(shutdown)
	}()

	return shutdown
}
