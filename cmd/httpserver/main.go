package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AliUnipal/chat/internal/bootstrapper"
	"github.com/AliUnipal/chat/internal/server/httpsrv/chatcontroller"
	_ "github.com/lib/pq"
)

type ctxKey string

const (
	srvCtx ctxKey = "serverContext"
)

type chatController interface {
	CreateChat(w http.ResponseWriter, r *http.Request)
	GetChats(w http.ResponseWriter, r *http.Request)
}

type userController interface {
	HandleCreateUser(w http.ResponseWriter, r *http.Request)
	HandleGetUser(w http.ResponseWriter, r *http.Request)
}

type messageController interface {
	HandleCreateMessage(w http.ResponseWriter, r *http.Request)
	HandleGetMessages(w http.ResponseWriter, r *http.Request)
}

type server struct {
	srv            *http.Server
	chatController chatController
	//userController    userController
	//messageController messageController
}

func NewServer(
	ctx context.Context,
	port int,
	chatController chatController,
	// userController userController,
	// controller messageController,
) *server {
	if port <= 0 || port > 65535 {
		panic(fmt.Sprintf("Invalid port number: %d", port))
	}

	if chatController == nil {
		panic("Chat controller cannot be nil")
	}
	//if userController == nil {
	//	panic("User controller cannot be nil")
	//}
	//if controller == nil {
	//	panic("Message controller cannot be nil")
	//}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.NewServeMux(),
		BaseContext: func(net.Listener) context.Context {
			return context.WithValue(ctx, srvCtx, "chatServer")
		},
	}

	s := &server{
		srv,
		chatController,
		//userController,
		//controller
	}
	s.registerHandlers()
	return s
}

func (s *server) registerHandlers() {
	mux, ok := s.srv.Handler.(*http.ServeMux)
	if !ok {
		panic("Handler is not of type *http.ServeMux")
	}

	mux.HandleFunc("POST /chats", s.chatController.CreateChat)
	mux.HandleFunc("GET /chats/{id}", s.chatController.GetChats)

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("pong"); err != nil {
			slog.ErrorContext(r.Context(), "Failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	bs, err := bootstrapper.New(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := bs.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	chatSvc := bs.NewChatService(ctx)
	chatController := chatcontroller.New(chatSvc)

	//userSvc := bs.NewUserService(ctx)
	//msgSvc := bs.NewMessageService(ctx)

	srv := NewServer(ctx, 8080, chatController)
	go func() {
		if err := srv.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("Server closed gracefully")
				return
			}

			slog.Error("Server failed to start", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		slog.Info("Received interrupt signal, shutting down server")
	case <-ctx.Done():
		slog.Info("Context cancelled, shutting down server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
}
