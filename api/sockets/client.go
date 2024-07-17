package sockets

import (
	"math/rand"
	"slices"
	"time"

	"github.com/aacebo/orca/api/sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	sockets sync.Map[string, []*Socket]
}

func New() Client {
	return Client{
		sockets: sync.NewMap[string, []*Socket](),
	}
}

func (self *Client) Get(appId string) *Socket {
	sockets := self.sockets.Get(appId)

	if sockets == nil || len(sockets) == 0 {
		return nil
	}

	seed := rand.New(rand.NewSource(time.Now().Unix()))
	return sockets[seed.Intn(len(sockets))]
}

func (self *Client) Add(appId string, conn *websocket.Conn) *Socket {
	sockets := self.sockets.Get(appId)

	if sockets == nil {
		sockets = []*Socket{}
	}

	socket := NewSocket(conn)
	sockets = append(sockets, socket)
	self.sockets.Set(appId, sockets)
	return socket
}

func (self *Client) Del(appId string, id string) {
	sockets := self.sockets.Get(appId)

	if sockets == nil {
		return
	}

	self.sockets.Set(appId, slices.DeleteFunc(sockets, func(s *Socket) bool {
		return s.ID == id
	}))
}
