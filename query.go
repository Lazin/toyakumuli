package main

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPServer instance
type HTTPServer struct {
	tss        *TimeSeriesStorage
	httpServer *http.Server
}

// NewHTTPServer creates new HTTPServer instance
func NewHTTPServer(port int, tss *TimeSeriesStorage) *HTTPServer {
	var srv HTTPServer
	srv.tss = tss
	srv.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        &srv,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &srv
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
