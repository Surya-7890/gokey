package server

import (
	"io"
	"net"
	"strings"
)

type Server struct {
	listener net.Listener
	address  string
	exit     chan struct{}
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		exit:    make(chan struct{}),
	}
}

func (s *Server) StartServer() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.listener = listener
	defer s.listener.Close()
	go s.AcceptConnections()

	<-s.exit

	return nil
}

func (s *Server) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.ReadFromConnections(conn)
	}
}

func (s *Server) ReadFromConnections(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 256)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				break
			}
			continue
		}
		data := buffer[:length]
		HanldeIncomingConnections(conn, strings.Split(strings.TrimSpace(string(data)), " "))
	}
}
