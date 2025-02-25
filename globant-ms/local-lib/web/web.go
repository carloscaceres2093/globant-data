package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	utils_log "globant-ms/local-lib/log"
	middleware_log "globant-ms/local-lib/web/middleware/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

func NewWebApp(environment string) *App {

	logger := utils_log.NewLogger(utils_log.WithLevel(utils_log.DebugLevel))
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

func newRouter(logger utils_log.Logger) *Router {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware_log.Logging(logger))

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

	w.Header().Set(ContentType, ApplicationJson)

	w.WriteHeader(code)

	if _, err := w.Write(jsonData); err != nil {
		log.Println(err.Error())
		return
	}
}
