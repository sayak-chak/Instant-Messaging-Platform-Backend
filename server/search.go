package server

import (
	"database/sql"
	"instant-messaging-platform-backend/config"

	"github.com/gofiber/fiber/v2"
)

var searchSuccessResponse string = `{"isUserPresent":true}`
var searchFailResponse string = `{"isUserPresent":false}`

func checkIfUserExists(ctx *fiber.Ctx) error {
	queryUsername := ctx.Query("username")
	if isAuthValid, err := IsAuthValid(string(ctx.Request().Header.Cookie("Auth-Token"))); !isAuthValid || err != nil {
		return fiber.NewError(400)
	}

	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlRes, err := db.Exec("select * from " + config.UsersTable + " where username='" + queryUsername + "'")

	noOfUsersWithGivenUsername, err := sqlRes.RowsAffected()
	if err != nil {
		return err
	}
	if noOfUsersWithGivenUsername != 1 {
		return ctx.Status(400).Send([]byte(searchFailResponse))
	}
	return ctx.Status(200).Send([]byte(searchSuccessResponse))

}
