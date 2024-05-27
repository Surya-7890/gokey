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
			delete(db, Map[key])
		}()
	}
}

func (p *Peer) GetData(key, database string) string {
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name"))
		return ""
	}
	return Map[key]
}

func (p *Peer) DeleteData(key, database string) error {
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name"))
		return fmt.Errorf("key not found")
	}
	delete(db, Map[key])
	return nil
}
