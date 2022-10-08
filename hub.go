package main

import (
	"errors"
	"github.com/gorilla/websocket"
)

var techConnection *Connection

type Connection struct {
	ws *websocket.Conn
}

func (c Connection) Listen(h *Hub) {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		h.Broadcast(msg)
	}
}

type Chatroom struct {
	client *Connection
	tech   *Connection
}

type Hub struct {
	chatrooms map[string]*Connection
}

func NewHub() *Hub {
	return &Hub{
		chatrooms: make(map[string]*Connection),
	}
}

func NewChatroom(guest *Connection) *Chatroom {
	if techConnection != nil {
		return nil
	}

	return &Chatroom{
		client: guest,
		tech:   techConnection,
	}
}

func (h *Hub) NewTechConnection(ws *websocket.Conn) error {
	if techConnection != nil {
		return errors.New("tech connection already exists")
	}

	techConnection = &Connection{
		ws: ws,
	}

	go techConnection.Listen(h)

	return nil
}
