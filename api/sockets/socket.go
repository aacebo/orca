package sockets

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Socket struct {
	ID        string
	CreatedAt time.Time

	conn *websocket.Conn
	log  *log.Logger
	mu   sync.RWMutex
}

func NewSocket(conn *websocket.Conn) *Socket {
	id := uuid.NewString()
	socket := Socket{
		ID:        id,
		CreatedAt: time.Now(),

		conn: conn,
		log:  log.New(os.Stdout, fmt.Sprintf("socket:%s ", id), log.Ldate|log.Ltime|log.Lshortfile),
		mu:   sync.RWMutex{},
	}

	go socket.onPing()
	return &socket
}

func (self *Socket) Read() (Event, error) {
	event := Event{}
	err := self.conn.ReadJSON(&event)

	if err != nil {
		self.log.Println(err.Error())
	}

	return event, err
}

func (self *Socket) Send(event Event) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	err := self.conn.WriteJSON(event)

	if err != nil {
		self.log.Println(err.Error())
		self.conn.Close()
	}

	return err
}

func (self *Socket) onPing() {
	for range time.Tick(20 * time.Second) {
		err := self.conn.WriteMessage(websocket.PingMessage, []byte{})

		if err != nil {
			self.log.Println(err.Error())
			self.conn.Close()
			return
		}
	}
}
