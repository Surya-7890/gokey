package server

import (
	"encoding/json"
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

func (p *Peer) SetData(key, val, database string, expiry time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()

	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name"))
		return
	}
	Map[key] = val
	if expiry != 0 {
		go func(key string) {
			time.Sleep(expiry)
			mutex.Lock()
			defer mutex.Unlock()
			delete(Map, key)
		}(key)
	}
	p.Conn.Write([]byte("success"))
}

func (p *Peer) GetData(key, database string) {
	mutex.Lock()
	defer mutex.Unlock()

	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	if key == "*" {
		data := make(map[string]string)
		for key, val := range Map {
			data[key] = val
		}
		data_map, err := json.Marshal(data)
		if err != nil {
			p.Conn.Write([]byte(err.Error()))
		}
		p.Conn.Write(data_map)
		return
	}
	p.Conn.Write([]byte(Map[key] + "\n"))
}

func (p *Peer) DeleteData(key, database string) {
	mutex.Lock()
	defer mutex.Unlock()

	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	delete(Map, key)
	p.Conn.Write([]byte("success"))
}

func (p *Peer) CreateTable(database string) {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := db[database]
	if ok {
		p.Conn.Write([]byte("database name already exists\n"))
		return
	}
	db[database] = make(map[string]string)
	p.Conn.Write([]byte("success"))
}
