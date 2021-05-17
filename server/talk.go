package server

import (
	"instant-messaging-platform-backend/client"

	"github.com/gofiber/websocket/v2"
)

var broadcastChannel = make(chan client.ChatMessage)

func talk(ctx *websocket.Conn) {
	defer ctx.Close() 

	go broadcastMsg()

	client.Write(ctx) 

	for { 
		var usernameAndMessage client.ChatMessage
		err := ctx.ReadJSON(&usernameAndMessage)
		if err != nil {
			client.Remove(ctx)
			break
		}
		broadcastChannel <- usernameAndMessage 
	}
}

func broadcastMsg() {
	for {
		msg := <-broadcastChannel
		for specificClient := range client.Clients {
			err := specificClient.WriteJSON(msg)
			if err != nil { // close that client
				specificClient.Close()
				client.Remove(specificClient)
			}
		}
	}
}
