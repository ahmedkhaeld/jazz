package jazz

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// routes gives the default routes that are available to the application if any
// and any middleware
func (j *Jazz) mux() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if j.Debug {
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)
	mux.Use(j.Trace)
	mux.Use(j.SessionLoad)
	mux.Use(j.NoSurf)

	return mux
}
