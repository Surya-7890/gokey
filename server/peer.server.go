package server

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Peer struct {
	Conn net.Conn
}

var (
	db    = make(map[string]map[string]string, 5)
	mutex = sync.Mutex{}
)

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Conn: conn,
	}
}

func (p *Peer) SetData(key, val, database string, expiry int) {
	mutex.Lock()
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name"))
		return
	}
	Map[key] = val
	mutex.Unlock()
	if expiry != 0 {
		go func() {
			time.Sleep(time.Duration(expiry) * time.Second)
			delete(Map, key)
		}()
	}
}

func (p *Peer) GetData(key, database string) {
	mutex.Lock()
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
	mutex.Unlock()
	p.Conn.Write([]byte(Map[key] + "\n"))
}

func (p *Peer) DeleteData(key, database string) {
	mutex.Lock()
	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	delete(Map, key)
	mutex.Lock()
	fmt.Println(Map[key])
}

func (p *Peer) CreateTable(database string) {
	mutex.Lock()
	_, ok := db[database]
	if ok {
		p.Conn.Write([]byte("database name already exists\n"))
		return
	}
	db[database] = make(map[string]string)
	mutex.Unlock()
	p.Conn.Write([]byte("table created\n"))
}
