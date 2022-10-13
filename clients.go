package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Id      string `json:"id"`
	TextMsg string `json:"textMsg"`
	DestId  string `json:"destId"`
}

func ServeTechWs(w http.ResponseWriter, r *http.Request, clients map[string]chan Message, Tech <-chan Message) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "text/json")
	var clientDisc = make(chan bool)
	defer close(clientDisc)

	go func() {
		defer func() { clientDisc <- true }()
		var message Message
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if msgType == websocket.CloseMessage {
				log.Println("Client disconnected")
				return
			}

			err = json.Unmarshal(msg, &message)
			id := message.Id
			conn.WriteJSON(message)
			clients[id] <- message
		}
	}()

	for {
		select {
		case msg := <-Tech:
			conn.WriteJSON(msg)
		case <-clientDisc:
			return
		}
	}

}

func ServeClientWs(w http.ResponseWriter, r *http.Request, tech chan<- Message, clients map[string]chan Message, id int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "text/json")
	defer conn.Close()
	var clientDisc = make(chan bool)
	client := clients[strconv.Itoa(id)]
	defer close(clientDisc)

	go func() {
		defer func() { clientDisc <- true }()
		var message Message
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if msgType == websocket.CloseMessage {
				log.Println("Closed connexion")
				return
			}

			err = json.Unmarshal(msg, &message)
			conn.WriteJSON(message)
			tech <- message
		}
	}()

	for {
		select {
		case msg := <-client:
			conn.WriteJSON(msg)
		case <-clientDisc:
			return
		}
	}

}
