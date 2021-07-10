package store

import (
	"instant-messaging-platform-backend/config"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var Sessions *session.Store

var sessionConfig = session.Config{
	Expiration:     30 * time.Minute,
	Storage:        nil,
	CookieName:     "Auth-Token",
	CookieDomain:   config.ClientDomain,
	CookiePath:     "",
	CookieSecure:   true,
	CookieHTTPOnly: false,
	CookieSameSite: "None",
}

func StartNewSession() {
	Sessions = session.New(sessionConfig)
}
