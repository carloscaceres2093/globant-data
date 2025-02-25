package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"globant-auth-ms/local-lib/log"
)

const (
	port = ":8080"
)

type Error struct {
	Message string `json:"message"`
}

type App struct {
	*Router
	Environment

	Log log.Logger

	server *http.Server
}
type Environment struct {
	Name string
}

func (wa *App) Run() error {
	wa.Log.Info(context.Background(), "", log.Field("port", port),
		log.Field("env", wa.Name))
	wa.server = &http.Server{Addr: port, Handler: wa.mux}

	shutdownComplete := handleShutdown(wa)

	err := wa.server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			wa.Log.Info(context.Background(), "", log.Field("error", err.Error()))
			return err
		}
		<-shutdownComplete
	}

	return nil
}

func (wa *App) Shutdown(ctx context.Context) error {
	wa.Log.Info(context.Background(), "", log.Field("message", port))

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

		time.Sleep(5 * time.Second)

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
