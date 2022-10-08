package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Ws(w http.ResponseWriter, r *http.Request, c chan []byte) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
	}

	go sendMsg(conn, c)
	go receiveMsg(conn, c)
}

func sendMsg(conn *websocket.Conn, c chan []byte) {
	defer conn.Close()
	for {
		conn.WriteMessage(websocket.TextMessage, <-c)
	}
}

func receiveMsg(conn *websocket.Conn, c chan []byte) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		c <- message
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}
