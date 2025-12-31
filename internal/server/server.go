package server

import (
	"net"
	"fmt"
)

type Server struct{
	listener net.Listener
	closed bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		closed: false,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed = true
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("can't close the listener %v\n", err)
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed {
				return //stops listening once closed
			}
			fmt.Printf("error listening: %v\n", err)
			s.closed = true
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" // Body
	conn.Write([]byte(response))
}