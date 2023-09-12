package ws

import (
	"log"
	"time"
	"wsserver/message"

	"github.com/gorilla/websocket"
)

func WriteData(ws *websocket.Conn, quit <-chan struct{}) {
	const tickerTime = 30 * time.Millisecond
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	msg, errMsg := message.Random()
	if errMsg != nil {
		return // TODO: Handle random struct error
	}

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
