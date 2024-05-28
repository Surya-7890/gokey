package main

import (
	"net"
	"strconv"
	"strings"
)

var Peers = make(map[net.Conn]*Peer, 10)

func HanldeIncomingConnections(conn net.Conn, message []string) {
	// parse incoming messages from connections
	switch strings.ToUpper(message[0]) {
	case "CREATE":
		// create a new table
		// format: CREATE tablename
		if len(message) < 2 {
			conn.Write([]byte("invalid number of arguments\n"))
			return
		}
		_, ok := Peers[conn]
		if !ok {
			Peers[conn] = &Peer{
				Conn: conn,
			}
		}
		Peers[conn].CreateTable(message[1])
	case "SET":
		// set key-value pair
		// format: SET key value tablename
		if len(message) < 4 {
			conn.Write([]byte("invalid number of arguments\n"))
			return
		}
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
		if len(message) < 3 {
			conn.Write([]byte("invalid number of arguments\n"))
			return
		}
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
		if len(message) < 5 {
			conn.Write([]byte("invalid number of arguments\n"))
			return
		}
		_, ok := Peers[conn]
		if !ok {
			Peers[conn] = &Peer{
				Conn: conn,
			}
		}
		expiry, err := strconv.Atoi(message[4])
		if err != nil {
			conn.Write([]byte("Invalid Time In Seconds\n"))
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
	default:
		conn.Write([]byte("invalid options, please try again\n"))
	}
}
