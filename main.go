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
	data := Message{}
loopRead:
	for {
		select {
		case <-quit:
			log.Println("received stop signal, stopping write")
			break loopRead
		default:
			err := ws.WriteJSON(data)
			if err != nil {
				log.Println("error writing message:", err)
				break loopRead
			}
			data.Data++
			time.Sleep(time.Millisecond)
		}
	}
}

func readMessage(ws *websocket.Conn, quit chan struct{}) {
loopRead:
	for {
		select {
		case <-quit:
			log.Println("received stop signal, stopping read")
			break loopRead
		default:
			msgType, _, err := ws.NextReader()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("received close message, terminating connection")
					close(quit)
					break loopRead
				}
				log.Println("error reading next message:", err)
				close(quit)
				break loopRead
			}

			if msgType == websocket.CloseMessage {
				log.Println("received close message, terminating connection")
				close(quit)
				break loopRead
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
	const (
		readSeconds  = 5
		writeSeconds = 10
		idleSeconds  = 15
	)
	http.HandleFunc("/", handler)
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  readSeconds * time.Second,
		WriteTimeout: writeSeconds * time.Second,
		IdleTimeout:  idleSeconds * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
