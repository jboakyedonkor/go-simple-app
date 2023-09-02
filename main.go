package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/caarlos0/env/v9"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

type config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func main() {
	cfg := config{}

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(jsonLogHandler)

	routes := Route{
		logger: logger,
	}

	tp := newTracerProvider()
	
	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	otel.SetTracerProvider(tp)

	router := chi.NewRouter()
	router.Use(routes.loggingMiddleware)
	router.Get("/host", routes.getHostname)
	router.Get("/health", routes.healthcheck)
	logger.Info("starting server listening on port " + strconv.Itoa(cfg.Port))
	server := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
