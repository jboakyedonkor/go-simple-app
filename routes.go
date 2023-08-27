package main

import (
	"encoding/json"
	"net/http"
	"os"

	"log/slog"
)

type Route struct {
	logger *slog.Logger
}
type HostnameResponse struct {
	Hostname string `json:"hostname"`
}

func (rt *Route) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		rt.logger.InfoContext(r.Context(), r.Method+" request", slog.String("path", r.URL.Path))
	})
}

func (rt *Route) getHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		rt.logger.ErrorContext(r.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(HostnameResponse{
		Hostname: hostname,
	})

	if err != nil {
		rt.logger.ErrorContext(r.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (rt *Route) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "up",
	})

	if err != nil {
		rt.logger.ErrorContext(r.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
