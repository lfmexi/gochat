package server

import (
	"log"
	"net"
)

// Handler is the interface that allows the manipulation of incoming TCP connections.
//
// Serves TCP connections to the implementation, giving the ability to handle the
// net.Conn typed variable that cames in its argument.
// Returns an error when something wrong happens to the connection.
type Handler interface {
	ServeTCP(net.Conn) error
}

// Server creates a TCP server and handles the incoming connections with
// an implementation of the Handler interface.
//
// Addr is the string address in which the server will listen to clients.
// Handler is the implementation of the Handler interface.
type Server struct {
	Addr     string
	Handler  Handler
	exitChan <-chan bool
}

// SetHandler takes an implementation of the Handler interface and sets it
// to s.
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

// Listen makes s to listen new tcp connections on s.Addr.
func (s *Server) Listen() {
	go s.listen()
	<-s.exitChan
}

func (s *Server) listen() {
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
			log.Printf("server accepting connection from %s", conn.RemoteAddr())
			go func() {
				log.Println("handling connection")
				s.handleConnection(conn)
			}()
		}
	}
}

// NewServer creates a new value for Server
// returns the pointer of that value.
func NewServer(address string, exitChan <-chan bool) *Server {
	return &Server{Addr: address, Handler: nil, exitChan: exitChan}
}
