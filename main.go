package main

import (
	"log"
	"net/http"
	"time"

	"wsserver/ws"

	"github.com/gorilla/websocket"
)

type server struct {
	upgrader websocket.Upgrader
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	quit := make(chan struct{})
	go ws.WriteData(conn, quit)
	go ws.ReadMessage(conn, quit)

	log.Println("connection alive:", r.RemoteAddr)
	<-quit
	log.Println("connection closed:", r.RemoteAddr)
}

func main() {
	const (
		readTime       = 5 * time.Second
		readHeaderTime = time.Second
		writeTime      = 10 * time.Second
		idleTime       = 15 * time.Second
	)
	server := &server{upgrader: websocket.Upgrader{}}
	http.HandleFunc("/", server.handler)
	srv := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       readTime,
		ReadHeaderTimeout: readHeaderTime,
		WriteTimeout:      writeTime,
		IdleTimeout:       idleTime,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
