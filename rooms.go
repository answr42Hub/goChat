package main

import (
	"math/rand"
	"time"
)

type Room struct {
	id         int32
	clients    map[*Clients]bool
	broadcast  chan *Message
	register   chan *Clients
	unregister chan *Clients
}

type Message struct {
	Message  string
	Type     string
	ClientID string
}

func NewRoom() *Room {
	rand.Seed(time.Now().UnixNano())
	room := &Room{
		id:         rand.Int31(),
		broadcast:  make(chan *Message),
		register:   make(chan *Clients),
		unregister: make(chan *Clients),
		clients:    make(map[*Clients]bool),
	}

	go room.run()

	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
