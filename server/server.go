package server

import (
	"log"
	"net"
)

type Handler interface {
	ServeTCP(net.Conn) error
}

type Server struct {
	Addr    string
	Handler Handler
}

func (s *Server) SetHandler(newHandler Handler) {
	s.Handler = newHandler
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		err := s.Handler.ServeTCP(conn)
		if err != nil {
			conn.Close()
			return
		}
	}
}

func (s *Server) Listen() {
	listen, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Server listening at %s\n", s.Addr)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		if s.Handler != nil {
			go s.handleConnection(conn)
		}
	}
}

func NewServer(address string) *Server {
	server := &Server{Addr: address, Handler: nil}
	return server
}
