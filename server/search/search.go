package search

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/server/authentication"

	"github.com/gofiber/fiber/v2"
)

var AuthValidator = authentication.IsAuthValid
var searchSuccessResponse string = `{"isUserPresent":true}`
var searchFailResponse string = `{"isUserPresent":false}`

func CheckIfUserExists(ctx *fiber.Ctx) error {
	queryUsername := ctx.Query("username")
	if isAuthValid, err := AuthValidator(string(ctx.Request().Header.Cookie("Auth-Token"))); !isAuthValid || err != nil {
		fmt.Println("CHANG(a)EEED ", queryUsername)
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
