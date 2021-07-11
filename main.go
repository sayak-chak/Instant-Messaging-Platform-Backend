package main

import (
	"fmt"
	"instant-messaging-platform-backend/database"
	"instant-messaging-platform-backend/server"
	"log"
)

func main() {
	err := database.SetupDataBase()
	if err != nil {
		panic("database initialization failed")
	}

	if err := server.New().Listen(":3000"); err != nil {
		log.Fatal("Error while running", err)
	}

	err = database.TearDown()
	if err != nil {
		fmt.Errorf("database teardown failed")
	}
}
