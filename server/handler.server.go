package server

import (
	"net"
	"strconv"
	"strings"
)

var Peers = make(map[net.Conn]*Peer, 10)

func HanldeIncomingConnections(conn net.Conn, message []string) {
	// parse incoming messages from connections
	switch strings.ToUpper(message[0]) {
	case "SET":
		// set key-value pair
		// format: SET key value tablename
		_, ok := Peers[conn]
		if !ok {
			Peers[conn] = &Peer{
				Conn: conn,
			}
		}
		Peers[conn].SetData(message[1], message[2], message[3], 0)
	case "GET":
		// get value by pair
		// format: GET key tablename
		_, ok := Peers[conn]
		if !ok {
			Peers[conn] = &Peer{
				Conn: conn,
			}
		}
		Peers[conn].GetData(message[1], message[2])
	case "SETEX":
		// set key-value pair with expiry in milliseconds
		// format: SET key value tablename expiry
		_, ok := Peers[conn]
		if !ok {
			Peers[conn] = &Peer{
				Conn: conn,
			}
		}
		expiry, err := strconv.Atoi(message[4])
		if err != nil {
			conn.Write([]byte("Invalid Time In Seconds"))
			return
		}
		Peers[conn].SetData(message[1], message[2], message[3], expiry)
	case "DELETE":
		// delete a pair using the key
		// format: DELETE key tablename
		_, ok := Peers[conn]
		if !ok {
			Peers[conn] = &Peer{
				Conn: conn,
			}
		}
		Peers[conn].DeleteData(message[1], message[2])
	}
}
