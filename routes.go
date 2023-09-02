package main

import (
	"encoding/json"
	"net/http"
	"os"

	"log/slog"

	"go.opentelemetry.io/otel"
)

type Route struct {
	logger *slog.Logger
}
type HostnameResponse struct {
	Hostname string `json:"hostname"`
}

func (rt *Route) loggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx, span := otel.Tracer("routes").Start(r.Context(), "Log")

		defer span.End()

		next.ServeHTTP(w, r.WithContext(newCtx))

		rt.logger.InfoContext(r.Context(), r.Method+" request", slog.String("path", r.URL.Path))
	})
}

func (rt *Route) getHostname(w http.ResponseWriter, r *http.Request) {

	_, span := otel.Tracer("routes").Start(r.Context(), "Host")

	defer span.End()

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

	_, span := otel.Tracer("routes").Start(r.Context(), "Healthcheck")

	defer span.End()
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
