package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func AddRoutesTo(app *fiber.App) {
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Get("/logincheck", LoginCheck)
	app.Get("/chat", websocket.New(chat))
	app.Get("/search", checkIfUserExists)
	app.Get("/history", getChatHistory)
}
