package server

import (
	"instant-messaging-platform-backend/server/store"

	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	app := fiber.New()
	store.StartNewSession()
	AddRoutesTo(app)

	return app
}
