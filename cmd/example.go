package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ctxKey string

const (
	srvCtx ctxKey = "serverContext"
)

func NewPingService() *pingService {
	return &pingService{}
}

type pingService struct{}

func (s *pingService) Ping() string {
	return "pong"
}

type Pinger interface {
	Ping() string
}

type dummyController struct {
	canceller context.CancelFunc
}

func (dc *dummyController) HandleQuit(w http.ResponseWriter, r *http.Request) {
	if dc.canceller != nil {
		dc.canceller()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"message": "Server is shutting down"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode quit response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(r.Context(), "Quit response sent", "response", resp)
}

func (dc *dummyController) HandleDummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"message": "dummy response"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode dummy response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(r.Context(), "Dummy response sent", "response", resp)
}

func (dc *dummyController) HandleLongRunning(w http.ResponseWriter, r *http.Request) {
	// Simulate a long-running operation
	select {
	case <-r.Context().Done():
		slog.InfoContext(r.Context(), "Long-running operation cancelled")
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	case <-time.After(5 * time.Second):
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"message": "long-running operation completed"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.ErrorContext(r.Context(), "Failed to encode long-running response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.InfoContext(r.Context(), "Long-running response sent", "response", resp)
}

func NewServer(ctx context.Context, port int, pingSvc Pinger, dummy *dummyController) *server {
	if port <= 0 || port > 65535 {
		panic(fmt.Sprintf("Invalid port number: %d", port))
	}

	if pingSvc == nil {
		panic("Ping service cannot be nil")
	}

	if dummy == nil {
		dummy = &dummyController{}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.NewServeMux(),
		BaseContext: func(net.Listener) context.Context {
			return context.WithValue(ctx, srvCtx, "pingServer")
		},
	}

	s := &server{srv, pingSvc, dummy}

	s.registerHandlers()
	return s
}

type server struct {
	srv     *http.Server
	pingSvc Pinger
	dc      *dummyController
}

func (s *server) registerHandlers() {
	mux, ok := s.srv.Handler.(*http.ServeMux)
	if !ok {
		panic("Handler is not of type *http.ServeMux")
	}

	mux.HandleFunc("/dummy", s.dc.HandleDummy)
	mux.HandleFunc("/long", s.dc.HandleLongRunning)
	mux.HandleFunc("/quit", s.dc.HandleQuit)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if name := r.URL.Query().Get("name"); name != "" {
			slog.InfoContext(r.Context(), "Received ping request with name", "name", name)
		}
		pong := s.pingSvc.Ping()
		resp := map[string]string{"message": pong}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		slog.InfoContext(r.Context(), "Ping response sent", "response", resp)
	})

	mux.HandleFunc("/ping/{name}", func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Received ping request", "name", r.PathValue("name"))
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})
}

func (s *server) Start() error {
	slog.Info("Starting server", "address", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server", "address", s.srv.Addr)
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	slog.Info("Server shutdown gracefully")
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Bootstrapper code:
	// bs := bootstrapper.New()
	// pingSvc := bs.NewPingService()

	pingSvc := NewPingService()
	dummyCtrl := &dummyController{canceller: cancel}

	srv := NewServer(ctx, 8080, pingSvc, dummyCtrl)
	go func() {
		if err := srv.Start(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("Server closed gracefully")
				return
			}

			slog.Error("Server failed to start", "error", err)
		}
	}()

	// Signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		slog.Info("Received interrupt signal, shutting down server")
	case <-ctx.Done():
		slog.Info("Context cancelled, shutting down server")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
}
