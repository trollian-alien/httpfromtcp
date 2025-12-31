package server

import (
	"fmt"
	"net"

	"github.com/trollian-alien/httpfromtcp/internal/response"
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
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}