package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"time"

	"github.com/AliUnipal/chat/internal/bootstrapper"
	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/server/httpsrv/chatcontroller"
	"github.com/AliUnipal/chat/internal/server/httpsrv/messagecontroller"
	"github.com/AliUnipal/chat/internal/server/httpsrv/usercontroller"
	"github.com/AliUnipal/chat/internal/server/wssrv/messagemanager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ctxKey string

const (
	srvCtxKey    ctxKey = "serverContext"
	userIdCtxKey ctxKey = "userIdContext"
)

type jwtManager interface {
	CreateToken(ctx context.Context, u user.User) (string, error)
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
}

type chatController interface {
	CreateChat(w http.ResponseWriter, r *http.Request)
	GetChats(w http.ResponseWriter, r *http.Request)
}

type userController interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserWithID(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type messageController interface {
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
}

type messageWsManager interface {
	ServeWS(w http.ResponseWriter, r *http.Request)
}

type contextAccessor struct {
}

func (ca *contextAccessor) CurrentUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIdCtxKey).(uuid.UUID)
	return id, ok
}

type server struct {
	srv               *http.Server
	mux               *http.ServeMux
	sess              jwtManager
	chatController    chatController
	userController    userController
	messageController messageController
	messageWsManager  messageWsManager
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.InfoContext(ctx, "Request received", "method", r.Method, "path", r.URL.Path)
		rec := httptest.NewRecorder()

		next.ServeHTTP(rec, r)
		slog.InfoContext(ctx, "Request processed", "status", rec.Code, "body", rec.Body.String())
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	})
}

func NewServer(
	ctx context.Context,
	port int,
	sess jwtManager,
	chatController chatController,
	userController userController,
	messageController messageController,
	messageWsManager messageWsManager,
) *server {
	if port <= 0 || port > 65535 {
		panic(fmt.Sprintf("Invalid port number: %d", port))
	}

	if sess == nil {
		panic("Session manager cannot be nil")
	}
	if chatController == nil {
		panic("Chat controller cannot be nil")
	}
	if userController == nil {
		panic("User controller cannot be nil")
	}
	if messageController == nil {
		panic("Type controller cannot be nil")
	}
	if messageWsManager == nil {
		panic("Message ws manager cannot be nil")
	}

	mux := http.NewServeMux()
	handler := With(
		mux,
		corsMiddleware,
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
		BaseContext: func(net.Listener) context.Context {
			return context.WithValue(ctx, srvCtxKey, "chatServer")
		},
	}

	s := &server{
		srv,
		mux,
		sess,
		chatController,
		userController,
		messageController,
		messageWsManager,
	}
	s.registerHandlers()
	return s
}

func corsMiddleware(next http.Handler) http.Handler {
	allowed := map[string]bool{
		"http://localhost:3000": true,
		"http://127.0.0.1:3000": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Add("Vary", "Access-Control-Request-Method")

			// Echo requested headers (most reliable)
			if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
				w.Header().Add("Vary", "Access-Control-Request-Headers")
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := r.Cookie("auth_token")
		if err != nil {
			slog.ErrorContext(ctx, "failed to get cookie", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		id, err := s.sess.VerifyToken(ctx, token.Value)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get session", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, userIdCtxKey, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *server) registerHandlers() {
	//mux, ok := s.srv.Handler.(*http.ServeMux)
	//if !ok {
	//	panic("Handler is not of type *http.ServeMux")
	//}
	mux := s.mux

	mux.HandleFunc("POST /chats", s.authenticated(s.chatController.CreateChat))
	mux.HandleFunc("GET /chats", s.authenticated(s.chatController.GetChats))

	mux.HandleFunc("POST /users", s.userController.CreateUser)
	mux.HandleFunc("GET /users/{id}", s.authenticated(s.userController.GetUserWithID))
	mux.HandleFunc("GET /identity", s.authenticated(s.userController.GetUser))
	mux.HandleFunc("POST /login", s.userController.Login)

	mux.HandleFunc("POST /chats/{id}/messages", s.authenticated(s.messageController.CreateMessage))
	mux.HandleFunc("GET /chats/{id}/messages", s.authenticated(s.messageController.GetMessages))

	// NOTE: Might not be the correct place but whatever!!!
	mux.HandleFunc("GET /ws/chats/{id}", s.authenticated(s.messageWsManager.ServeWS))

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
	//if err := s.srv.ListenAndServeTLS("cert/localhost+2.pem", "cert/localhost+2-key.pem"); err != nil && err != http.ErrServerClosed {
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
	privKey, pubKey, err := loadJWTKeys()
	if err != nil {
		log.Fatal(err)
	}

	bs, err := bootstrapper.New(dbUrl, privKey, pubKey)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := bs.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	chatSvc := bs.NewChatService(ctx)
	chatController := chatcontroller.New(chatSvc, &contextAccessor{})

	userSvc := bs.NewUserService(ctx)
	userController := usercontroller.New(userSvc, &contextAccessor{})

	msgSvc := bs.NewMessageService(ctx)
	msgController := messagecontroller.New(msgSvc)
	msgWsManager := messagemanager.New(msgSvc, &contextAccessor{})

	srv := NewServer(
		ctx, 8080,
		bs.NewSessionManager(ctx),
		chatController,
		userController,
		msgController,
		msgWsManager,
	)
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

// With applies the given middlewares to the given handler in the order they
// are given. That means that the first middleware is the outermost one.
func With(
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {
	// Do this in reverse order so that the first middleware is the outermost one
	if len(middlewares) == 0 {
		return handler
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		if m := middlewares[i]; m != nil {
			if h := m(handler); h != nil {
				handler = h
			}
		}
	}

	return handler
}

func WithHandlerFunc(
	handler http.HandlerFunc,
	middlewares ...func(http.Handler) http.Handler,
) http.HandlerFunc {
	return With(handler, middlewares...).ServeHTTP
}

func (s *server) authenticated(h http.HandlerFunc) http.HandlerFunc {
	return WithHandlerFunc(h, s.authMiddleware)
}

func loadJWTKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKeyData, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	pubKeyData, err := os.ReadFile(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return &rsa.PrivateKey{}, &rsa.PublicKey{}, err
	}

	return privKey, pubKey, nil
}
