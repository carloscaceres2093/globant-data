package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	utilsLog "globant-api/local-lib/log"
	middlewareLog "globant-api/local-lib/web/middleware/log"
	middlewareAuth "globant-api/local-lib/web/middleware/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

func NewWebApp(environment string) *App {

	logger := utilsLog.NewLogger(utilsLog.WithLevel(utilsLog.DebugLevel))
	router := newRouter(logger)
	env := newEnvironment(environment)

	return &App{
		Router:      router,
		Environment: env,
		Log:         logger,
	}
}

func newEnvironment(environment string) Environment {
	return Environment{
		Name: environment,
	}
}

func newRouter(logger utilsLog.Logger) *Router {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middlewareLog.Logging(logger))
	mux.Use(middlewareAuth.AuthMiddleware)
	router := &Router{
		mux: mux,
	}
	return router
}


func RespondJSON(w http.ResponseWriter, v interface{}, code int) {
	
	if code == http.StatusNoContent || v == nil {
		w.WriteHeader(code)
		return
	}

	var jsonData []byte

	var err error
	switch v := v.(type) {
	case []byte:
		jsonData = v
	case io.Reader:
		jsonData, err = io.ReadAll(v)
	default:
		jsonData, err = json.Marshal(v)
	}

	if err != nil {
		log.Println(err.Error())
		return
	}

	// Set the content type.
	w.Header().Set(ContentType, ApplicationJson)

	// Write the status code to the response and context.
	w.WriteHeader(code)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		log.Println(err.Error())
		return
	}
}
