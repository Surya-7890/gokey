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

type Message struct {
	val       string
	expiry    *time.Duration
	createdAt *time.Time
}

var (
	db    = make(map[string]map[string]*Message, 5)
	mutex = sync.Mutex{}
)

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		Conn: conn,
	}
}

func (p *Peer) SetData(key, val, database string) {
	mutex.Lock()
	defer mutex.Unlock()

	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	Map[key] = &Message{
		val: val,
	}
	p.Conn.Write([]byte("success\n"))
}

func (p *Peer) SetDataWithExpiration(key, val, database string, expiry time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()

	Map, ok := db[database]
	if !ok {
		p.Conn.Write([]byte("invalid database name\n"))
		return
	}
	now := time.Now()
	Map[key] = &Message{
		val:       val,
		expiry:    &expiry,
		createdAt: &now,
	}
	p.Conn.Write([]byte("success\n"))
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
			if val.expiry != nil {
				duration := time.Since(*val.createdAt)
				if duration > *val.expiry {
					delete(Map, key)
					continue
				}
			}
			data[key] = val.val
		}
		data_map, err := json.Marshal(data)
		if err != nil {
			p.Conn.Write([]byte(err.Error() + "\n"))
		}
		p.Conn.Write([]byte(data_map))
		p.Conn.Write([]byte("\n"))
		return
	}
	if Map[key].expiry != nil {
		if time.Since(*Map[key].createdAt) < *Map[key].expiry {
			p.Conn.Write([]byte(Map[key].val + "\n"))
			return
		}
	}
	delete(Map, key)
	p.Conn.Write([]byte("\n"))
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
	p.Conn.Write([]byte("success\n"))
}

func (p *Peer) CreateTable(database string) {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := db[database]
	if ok {
		p.Conn.Write([]byte("database name already exists\n"))
		return
	}
	db[database] = make(map[string]*Message)
	p.Conn.Write([]byte("success\n"))
}
