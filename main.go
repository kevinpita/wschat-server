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

type server struct {
	upgrader websocket.Upgrader
}

func writeData(ws *websocket.Conn, quit <-chan struct{}) {
	msg := Message{}

	const tickerTime = 30 * time.Millisecond
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			log.Println("received stop signal, stopping write")
			return
		case <-ticker.C:
			err := ws.WriteJSON(msg)
			if err != nil {
				log.Println("error writing message:", err)
				return
			}
			msg.Data++
		}
	}
}

func readMessage(ws *websocket.Conn, quit chan<- struct{}) {
	for {
		msgType, msgData, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("connection has been closed by the client")
			} else {
				log.Println("error reading next message:", err)
			}
			close(quit)
			return
		}

		if msgType == websocket.CloseMessage {
			log.Println("received close message, terminating connection")
			close(quit)
			return
		}

		log.Println("received data:", msgData)
	}
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	quit := make(chan struct{})
	go writeData(conn, quit)
	go readMessage(conn, quit)

	log.Println("connection alive:", r.RemoteAddr)
	<-quit
	log.Println("connection closed:", r.RemoteAddr)
}

func main() {
	const (
		readTime  = 5 * time.Second
		writeTime = 10 * time.Second
		idleTime  = 15 * time.Second
	)
	server := &server{upgrader: websocket.Upgrader{}}
	http.HandleFunc("/", server.handler)
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  readTime,
		WriteTimeout: writeTime,
		IdleTimeout:  idleTime,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
