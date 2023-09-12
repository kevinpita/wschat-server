package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

func ReadMessage(ws *websocket.Conn, quit chan<- struct{}) {
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
