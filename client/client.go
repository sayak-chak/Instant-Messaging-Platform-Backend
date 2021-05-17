package client

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

var Clients = make(map[*websocket.Conn]bool) //indicates which connections are open, remove if client closes
var lock = sync.RWMutex{}

type ChatMessage struct{
	Username string `json:"username"`
	Msg  string `json:"message"`

}

func Write(ctx *websocket.Conn) {
	lock.Lock()
	defer lock.Unlock()
	Clients[ctx] = true
}

func Remove(ctx *websocket.Conn) {
	lock.Lock()
	defer lock.Unlock()
	delete(Clients, ctx)
}
