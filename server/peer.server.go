package server

import (
	"fmt"
	"net"
	"time"
)

type Peer struct {
	Conn net.Conn
}

var db = make(map[string]map[string]string, 5)

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Conn: conn,
	}
}

func (p *Peer) SetData(key, val, database string, expiry int) {
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name"))
		return
	}
	Map[key] = val
	if expiry != 0 {
		go func() {
			time.Sleep(time.Duration(expiry) * time.Second)
			delete(Map, key)
		}()
	}
}

func (p *Peer) GetData(key, database string) {
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	if key == "*" {
		for key, val := range Map {
			data := fmt.Sprintf("%s : %s\n", key, val)
			p.Conn.Write([]byte(data))
		}
		return
	}
	p.Conn.Write([]byte(Map[key] + "\n"))
}

func (p *Peer) DeleteData(key, database string) {
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	delete(Map, key)
	fmt.Println(Map[key])
}

func (p *Peer) CreateTable(database string) {
	_, ok := db[database]
	if ok {
		p.Conn.Write([]byte("database name already exists\n"))
		return
	}
	db[database] = make(map[string]string)
	p.Conn.Write([]byte("table created\n"))
}
