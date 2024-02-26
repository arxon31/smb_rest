package httpserver

import (
	"fmt"
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func Addr(host, port string) Option {
	return func(s *Server) {
		s.server.Addr = fmt.Sprintf("%s:%s", host, port)
	}
}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
