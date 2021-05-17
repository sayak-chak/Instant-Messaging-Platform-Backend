package main

import (
	"instant-messaging-platform-backend/server"
	"log"
)

func main() {
	if err := server.New().Listen(":3000"); err != nil{
		log.Fatal("Error while running", err)
	}
}