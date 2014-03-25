package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

type Server struct {
	Address string

	connLimit int
	conns     chan struct{}

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

func (s *Server) LimitConnections(count int) {
	s.connLimit = count
	s.conns = make(chan struct{}, count)
	for i := 0; i < count; i++ {
		s.conns <- struct{}{}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.connLimit != 0 {
		select {
		case <-s.conns:
			defer func() {
				s.conns <- struct{}{}
			}()
			break
		default:
			log.Printf("Connection from %s refused. Limit reached.", r.RemoteAddr)
			http.Error(w, "Connection limit reached", http.StatusServiceUnavailable)
			return
		}
	}

	conn, err := websocket.Upgrade(w, r, nil, ReadBufferSize, WriteBufferSize)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Websocket handshake expected", 400)
		return
	} else if err != nil {
		log.Print(err)
		return
	}

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
