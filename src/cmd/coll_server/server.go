package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

type Server struct {
	Address   string
	ConnLimit int

	conns int
	m     sync.Mutex

	Http *http.Server
}

func NewServer(address string) *Server {
	s := &Server{
		Address: address,
		Http: &http.Server{
			Addr: address,
		},
	}
	s.Http.Handler = s
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.ConnLimit != 0 {
		s.m.Lock()
		if s.conns >= s.ConnLimit {
			s.m.Unlock()
			log.Printf("Connection from %s refused. Limit reached.", r.RemoteAddr)
			http.Error(w, "Connection limit reached", http.StatusServiceUnavailable)
			return
		}
		s.m.Unlock()
	}

	conn, err := websocket.Upgrade(w, r, nil, ReadBufferSize, WriteBufferSize)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Websocket handshake expected", 400)
		return
	} else if err != nil {
		log.Print(err)
		return
	}

	s.m.Lock()
	s.conns++
	s.m.Unlock()
	defer func() {
		s.m.Lock()
		s.conns--
		s.m.Unlock()
	}()

	log.Printf("Websocket connection from %s established.", r.RemoteAddr)
	// serve websocket connection
	for {
		msgType, _, err := conn.NextReader()
		if err != nil {
			log.Printf("Websocket connection from %s closed.", r.RemoteAddr)
			return
		}
		switch msgType {
		case websocket.BinaryMessage:
		case websocket.TextMessage:
		}
	}
}
