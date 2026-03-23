package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

const Port = ":8080"

func NewServer(log *slog.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mdw := LogMiddleware(log)
	handler := mdw(mux)

	return &http.Server{
		Addr:    Port,
		Handler: handler,
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func LogMiddleware(log *slog.Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.InfoContext(r.Context(), fmt.Sprintf("request: %s %s", r.Method, r.URL.Path))
			next.ServeHTTP(w, r)
		})
	}
}
