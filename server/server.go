package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

var Sessions *session.Store

func New() *fiber.App {
	app := fiber.New()

	sessionConfig := session.Config{
		Expiration:     30 * time.Minute,
		Storage:        nil,
		CookieName:     "Auth-Token",
		CookieDomain:   "",
		CookiePath:     "",
		CookieSecure:   false,
		CookieHTTPOnly: false,
	}

	// create session handler
	Sessions = session.New(sessionConfig)

	AddRoutesTo(app)
	return app
}
