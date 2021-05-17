package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

)

func addRoutesTo(app *fiber.App) { 
	app.Get("/register", websocket.New(register)) //TODO: change to restful architecture

	app.Get("/login", websocket.New(login))

	app.Get("/talk", websocket.New(talk))
}
