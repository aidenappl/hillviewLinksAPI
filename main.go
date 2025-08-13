package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hillview.tv/linksAPI/db"
	"github.com/hillview.tv/linksAPI/env"
	"github.com/hillview.tv/linksAPI/middleware"
	"github.com/hillview.tv/linksAPI/routers"
)

func main() {
	// Ping DB
	if err := db.PingDB(); err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	} else {
		log.Println("✅ Connected to the database successfully")
	}

	primary := mux.NewRouter()

	// Healthcheck Endpoint

	primary.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	// Define the API Endpoints

	r := primary.PathPrefix("/links/v1.1").Subrouter()

	// Logging of requests
	r.Use(middleware.LoggingMiddleware)

	// Adding response headers
	r.Use(middleware.MuxHeaderMiddleware)

	// Track & Update Last Active
	r.Use(middleware.TokenHandlers)

	r.HandleFunc("/check/{route}", routers.CheckLinkRouteHandler).Methods(http.MethodGet)

	// Launch API Listener
	fmt.Printf("✅ Hillview Links API running on port %s\n", env.Port)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Authorization", "Accept", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(originsOk, headersOk, methodsOk)(primary)))
}
