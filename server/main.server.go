package server

import (
	"fmt"
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

func (s *Server) StartServer() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	defer s.listener.Close()
	go s.AcceptConnections()
	<-s.exit
}

func (s *Server) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("connection from %v failed\n", conn.RemoteAddr().String())
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
				fmt.Printf("connection %s closed\n", conn.RemoteAddr().String())
				conn.Close()
				break
			}
			fmt.Printf("message from %v could not be read\n", conn.RemoteAddr().String())
			continue
		}
		data := buffer[:length]
		HanldeIncomingConnections(conn, strings.Split(strings.TrimSpace(string(data)), " "))
	}
}
