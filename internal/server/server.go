package server

import (
	"fmt"
	"net"

	"github.com/trollian-alien/httpfromtcp/internal/request"
	"github.com/trollian-alien/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct{
	handler  Handler
	listener net.Listener
	closed bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		handler: handler,
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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}