package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Server ...
type Server struct {
	httpServer *http.Server

	upgrader *websocket.Upgrader

	// signals
	addNewClientSignal chan *websocket.Conn
}

// NewServer ...
func NewServer(opts ...option) *Server {
	emptyServer := &Server{
		httpServer: &http.Server{},
	}
	emptyServer.httpServer.Handler = http.HandlerFunc(emptyServer.AcceptNewClient)

	for _, opt := range opts {
		if err := opt(emptyServer); err != nil {
			panic(err)
		}
	}

	return emptyServer
}

func (s *Server) Start() error {
	go s.startServerLoop()
	return s.httpServer.ListenAndServe()
}

func (s *Server) AcceptNewClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: add graceful shutdown
		log.Println(err)
		return
	}

	s.addNewClientSignal <- conn
}

func (s *Server) startServerLoop() {
	for {
		conn := <-s.addNewClientSignal
		_ = conn
	}
}
