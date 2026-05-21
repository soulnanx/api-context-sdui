package server

import (
	"log"
	"net/http"
	"time"

	"api-context-sdui/internal/screen"
)

type Server struct {
	handler http.Handler
}

func NewServer() *Server {
	mux := http.NewServeMux()
	screenHandler := &screen.ScreenHandler{}

	mux.HandleFunc("GET /api/screen/{screenId}", screenHandler.GetScreen)

	handler := loggingMiddleware(mux)
	return &Server{handler: handler}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
