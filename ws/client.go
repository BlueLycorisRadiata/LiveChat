package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomid"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomid"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer c.conn.Close()
	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.conn.Close()
	}()
	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
			) {
				log.Println("client disconnected")
			} else {
				log.Printf("real error: %v", err)
			}
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}
		hub.Broadcast <- msg
	}
}
