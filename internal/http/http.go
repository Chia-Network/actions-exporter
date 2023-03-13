package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	internalProm "github.com/chia-network/actions-exporter/internal/prometheus"
)

// Server is the config for the http server
type Server struct {
	port uint16
}

// NewServer returns a new instance of the http server
func NewServer(port uint16) (*Server, error) {
	server := &Server{
		port: port,
	}

	return server, nil
}

// StartServer starts the metrics server
func (s *Server) StartServer() error {
	log.Printf("Starting server on port %d", s.port)

	// Create new router with healthcheck handler
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthcheck).Methods(http.MethodGet)

	// Unregister default prometheus collectors so we don't collect a bunch of pointless metrics
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	// Set metrics handler
	prometheus.MustRegister(internalProm.NewCollector())
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	return srv.ListenAndServe()
}

// healthcheck an unprotected endpoint that just reports an http 200 if the server is still responding to requests
func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Errorf("failed to write ok response for healthz endpoint: %s", err.Error())
	}
}
