package server

import "github.com/gorilla/websocket"

type option func(s *Server) error

// WithPort ...
func WithPort(httpPort string) option {
	return func(s *Server) error {
		s.httpServer.Addr = ":" + httpPort
		return nil
	}
}

// WithNewClientBufferSize ...
func WithNewClientBufferSize(size int) option {
	return func(s *Server) error {
		s.addNewClientSignal = make(chan *websocket.Conn, size)
		return nil
	}
}
