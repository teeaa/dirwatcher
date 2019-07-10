package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func serveFiles(w http.ResponseWriter, r *http.Request) {
	log.Info("Requested list of files")
	err := json.NewEncoder(w).Encode(map[string][]map[string]string{"files": formatFiledata()})
	if err != nil {
		log.Error("Error encoding JSON response: ", err)
	}
}

// Router definition in separate func for testing
func getRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	router.Use(setHeaders)
	router.HandleFunc("/files", serveFiles).Methods("GET")

	return router
}

// Start HTTP server in it's own thread and return pointer to it for shutdown
func startRestAPI() *http.Server {
	log.Info("Starting REST API")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      getRouter(),
	}

	go func() {
		// We might get ErrServerClosed when exiting the server, but that is alright
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Unable to start REST API: ", err)
		}
	}()

	return srv
}

// Graceful killing of server
func shutdownRestAPI(srv *http.Server) {
	log.Info("Shutting down REST API")
	err := srv.Close()
	if err != nil {
		log.Error("Error closing REST API: ", err)
	}
}

// We will only return JSON for our defined endpoint
// (Added for possible future consistency)
func setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
