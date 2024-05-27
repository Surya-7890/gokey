package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	listener net.Listener
	address  string
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
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
				conn.Close()
			}
			fmt.Printf("message from %v could not be read\n", conn.RemoteAddr().String())
			continue
		}
		data := buffer[:length]
		fmt.Println(string(data))
	}
}
