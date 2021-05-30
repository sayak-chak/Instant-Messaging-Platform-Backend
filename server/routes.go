package server

import (
	"instant-messaging-platform-backend/server/chat"
	"instant-messaging-platform-backend/server/entry"
	"instant-messaging-platform-backend/server/search"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func AddRoutesTo(app *fiber.App) {
	app.Post("/register", entry.Register)
	app.Post("/login", entry.Login)
	app.Get("/logincheck", entry.LoginCheck)
	app.Get("/chat", websocket.New(chat.Chat))
	app.Get("/search", search.CheckIfUserExists)
	app.Get("/history", chat.GetChatHistory)
}
