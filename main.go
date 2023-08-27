package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Data int `json:"data"`
}

func writeData(ws *websocket.Conn, quit <-chan struct{}) {
	data := Message{
		Data: 0,
	}
	for {
		select {
		case <-quit:
			log.Println("received stop signal, stopping write")
			return
		default:
			err := ws.WriteJSON(data)
			if err != nil {
				log.Println("error writing message:", err)
				break
			}
			data.Data++
		}
	}
}

func readMessage(ws *websocket.Conn, quit chan struct{}) {
	for {
		select {
		case <-quit:
			log.Println("received stop signal, stopping read")
			return
		default:
			msgType, _, err := ws.NextReader()
			if err != nil {
				log.Println("error reading next message:", err)
				close(quit)
				break
			}

			if msgType == websocket.CloseMessage {
				log.Println("received close message, terminating connection")
				close(quit)
				break
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	quit := make(chan struct{})
	go writeData(conn, quit)
	go readMessage(conn, quit)

	<-quit
	log.Println("closing connection:", r.RemoteAddr)
}

func main() {
	http.HandleFunc("/", handler)
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
