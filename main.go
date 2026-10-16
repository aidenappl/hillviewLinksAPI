package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hillview.tv/linksAPI/background"
	"github.com/hillview.tv/linksAPI/db"
	"github.com/hillview.tv/linksAPI/env"
	"github.com/hillview.tv/linksAPI/middleware"
	"github.com/hillview.tv/linksAPI/routers"
	"github.com/hillview.tv/linksAPI/telemetry"
)

func main() {
	telemetry.Init(telemetry.Config{
		Service:   "hillview-links-api",
		IngestURL: env.MonitorIngestURL,
		APIKey:    env.MonitorAPIKey,
		Env:       env.MonitorEnv,
		SpoolDir:  env.MonitorSpoolDir,
		Stdout:    env.MonitorStdout,
	})

	// Ping DB
	if err := db.PingDB(); err != nil {
		telemetry.Fatal("service.startup.db_unreachable", err, nil)
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	} else {
		log.Println("✅ Connected to the database successfully")
	}

	primary := mux.NewRouter()

	// Add context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Healthcheck Endpoint
	primary.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Define the API Endpoints

	r := primary.PathPrefix("/links/v1.1").Subrouter()

	// Monitor: one http.request.end event per request, panics recovered
	r.Use(telemetry.Middleware)

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// Adding response headers
	r.Use(middleware.MuxHeaderMiddleware)

	// Track & Update Last Active
	r.Use(middleware.TokenHandlers)

	r.HandleFunc("/check/{route}", routers.CheckLinkRouteHandler).Methods(http.MethodGet)

	// Launch API Listener
	fmt.Printf("✅ Hillview Links API running on port %s\n", env.Port)

	// Add the health check polling
	telemetry.Go(ctx, "background.healthcheck", background.StartHealthCheckPolling)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"https://team.hillview.tv", "https://hillview.tv", "https://live.hillview.tv", "http://localhost:3000", "http://localhost:3001", "http://localhost:3002"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	exposedOk := handlers.ExposedHeaders([]string{"X-Request-ID"})

	server := &http.Server{
		Addr:              ":" + env.Port,
		Handler:           handlers.CORS(originsOk, headersOk, exposedOk, methodsOk)(primary),
		ReadHeaderTimeout: 15 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	if err := telemetry.ServeWithGracefulShutdown(server, stop); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
