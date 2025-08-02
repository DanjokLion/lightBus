package http

import (
	"context"
	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/DanjokLion/lightBus/internal/broker"
	"github.com/DanjokLion/lightBus/pkg/logger"
)

type Server struct {
	bus 	broker.Broker
	log 	*logger.Logger
	server	*http.Server
}

func NewServer(bus broker.Broker, log *logger.Logger) *Server {
	mux := http.NewServeMux()

	s := &Server{
		bus: bus,
		log: log,
		server: &http.Server{
			Handler: mux,
		},
	}

	mux.HandleFunc("/publish", s.handlePublish)
	mux.HandleFunc("/dlq", s.handleDLQ)
	mux.HandleFunc("/healthz", s.handleHealth)

	return s
}

func (s *Server) Start(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Topic 	string 	`json:"topic"`
		Data 	string	`json:"data"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := s.bus.Publish(r.Context(), req.Topic, []byte(req.Data)); err != nil {
		s.log.Error("Publish error: %v", err)
		http.Error(w, "publish failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("published"))
}

func (s *Server) handleDLQ(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "topic is required", http.StatusBadRequest)
		return
	}

	dlq := s.bus.GetDLQ(topic)
	resp, _ := json.Marshal(dlq)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}