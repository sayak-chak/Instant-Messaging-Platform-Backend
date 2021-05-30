package server

import (
	"instant-messaging-platform-backend/database"
	"instant-messaging-platform-backend/server/store"

	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	err := database.SetupDataBase()
	if err != nil {
		panic("Database initialization failed")
	}

	app := fiber.New()

	store.StartNewSession()

	AddRoutesTo(app)
	return app
}
