package client

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

var ClientsMap = make(map[string]*websocket.Conn) //indicates which connections are open, remove if client closes
var lock = sync.RWMutex{}

func RemoveClient(username string) {
	lock.Lock()
	defer lock.Unlock()
	delete(ClientsMap, username)
}

func AddClient(username string, ctx *websocket.Conn) {
	if _, isPresent := ClientsMap[username]; isPresent {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	ClientsMap[username] = ctx
}
