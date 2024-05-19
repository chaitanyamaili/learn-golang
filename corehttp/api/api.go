package api

import (
	"log/slog"
	"net/http"
)

// AppServer is the main server struct
type AppServer struct {
	addr string
	log  *slog.Logger
}

// NewAppServer creates a new server instance
func NewAppServer(addr string, log *slog.Logger) *AppServer {
	return &AppServer{
		addr: addr,
		log:  log,
	}
}

// Run starts the server
func (s *AppServer) Run() error {
	router := http.NewServeMux()

	// Define the handler for the / route
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		uID := r.PathValue("userID")
		w.Write([]byte("User ID: " + uID))
	})

	// Create a new router for the /api/v1/ route
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	middlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		RequestAuthMiddleware,
	)

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(router, s.log),
	}

	s.log.Info("Starting server")

	return server.ListenAndServe()
}

// RequestLoggerMiddleware logs the incoming requests
func RequestLoggerMiddleware(next http.Handler, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Method " + r.Method + " requested on " + r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// RequestAuthMiddleware is a middleware to authenticate requests
func RequestAuthMiddleware(next http.Handler, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer my-secret-token" {
			log.Error("Unauthorized request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Middleware is a type for middleware functions
type Middleware func(http.Handler, *slog.Logger) http.HandlerFunc

// MiddlewareChain chains multiple middleware functions
func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler, log *slog.Logger) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next, log)
		}

		return next.ServeHTTP
	}
}
