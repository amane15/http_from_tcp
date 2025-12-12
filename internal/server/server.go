package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/amane15/http_from_tcp/internal/request"
	"github.com/amane15/http_from_tcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}

	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
		}
		go s.handle(conn)
	}
}

// func (s *Server) handle(conn net.Conn) {
// 	defer conn.Close()
//
// 	req, err := request.RequestFromReader(conn)
// 	if err != nil {
// 		hErr := &HandlerError{
// 			StatusCode: response.StatusCodeBadRequest,
// 			Message:    err.Error(),
// 		}
// 		hErr.Write(conn)
// 		return
// 	}
// 	buf := bytes.NewBuffer([]byte{})
// 	hErr := s.handler(buf, req)
// 	if hErr != nil {
// 		hErr.Write(conn)
// 		return
// 	}
// 	b := buf.Bytes()
//
// 	response.WriteStatusLine(conn, response.StatusCodeSuccess)
// 	headers := response.GetDefaultHeaders(len(b))
// 	response.WriteHeaders(conn, headers)
// 	conn.Write(b)
// }

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
