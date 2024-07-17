package agents

import (
	"net/http"

	"github.com/aacebo/orca/api/amqp"
	"github.com/aacebo/orca/api/common"
	"github.com/aacebo/orca/api/sockets"
	"github.com/google/uuid"

	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

func Handler(ctx common.Context) http.HandlerFunc {
	amqp := ctx.Value("amqp").(*amqp.Client)
	socks := ctx.Value("sockets").(*sockets.Client)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			render.Status(r, 500)
			render.PlainText(w, r, err.Error())
			return
		}

		socket := socks.Add(id, conn)

		defer func() {
			conn.Close()
			socks.Del(id, socket.ID)
		}()

		for {
			event, err := socket.Read()

			if err != nil {
				return
			}

			if event.Type == sockets.Ack {
				amqp.Ack(event.ID)
			}
		}
	}
}
