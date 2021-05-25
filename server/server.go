package server

import (
	// "fmt"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	// "github.com/gofiber/websocket/v2"
)

var Sessions *session.Store

func New() *fiber.App {
	err := database.SetupDataBase()
	if err != nil {
		panic("Database initialization failed")
	}

	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "",
	// 	AllowCredentials: true,
	// }))

	sessionConfig := session.Config{
		Expiration:     30 * time.Minute,
		Storage:        nil,
		CookieName:     "Auth-Token",
		CookieDomain:   config.ClientDomain,
		CookiePath:     "",
		CookieSecure:   true,
		CookieHTTPOnly: false,
		CookieSameSite: "None",

		// Expiration:     30 * time.Minute,
		// Storage:        nil,
		// CookieName:     "Auth-Token",
		// CookieDomain:   "",
		// CookiePath:     "",
		// CookieSecure:   false,
		// CookieHTTPOnly: false,
	}

	// create session handler
	Sessions = session.New(sessionConfig)

	AddRoutesTo(app)
	return app
}
