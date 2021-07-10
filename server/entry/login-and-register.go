package entry

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database/model"
	"instant-messaging-platform-backend/server/store"

	"encoding/json"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

type LoginCheckResponse struct {
	Username  string `json:"username"`
	AuthToken string `json:"auth"`
}
type LoginRegisterResponse struct {
	IsRequestSuccessful bool `json:"isRequestSuccessful"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var loginRegisterFailResponse string = `{"isRequestSuccessful":false}`

func Register(ctx *fiber.Ctx) error {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	loginReq := LoginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return err
	}

	err := sendSuccessResponse(insertCreds(db, loginReq.Username, loginReq.Password), ctx, 201, 409)
	if err != nil {
		return fiber.NewError(409, loginRegisterFailResponse)
	}

	return storeAuthCookieAndUpdateDatabase(ctx, db)
}

func Login(ctx *fiber.Ctx) error {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	loginReq := LoginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return err
	}

	err := sendSuccessResponse(validateCreds(db, loginReq.Username, loginReq.Password), ctx, 200, 400)
	if err != nil {
		return fiber.NewError(400, loginRegisterFailResponse)
	}

	return storeAuthCookieAndUpdateDatabase(ctx, db)
}

func LoginCheck(ctx *fiber.Ctx) error {
	//var uuid string
	var username string
	authToken := string(ctx.Request().Header.Cookie("Auth-Token"))

	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	err := db.Model(&model.CredsTable{}).Column("username").Where("uuid = ?", authToken).Select(&username)

	if err != nil {
		return err
	}
	loginCheckResp, err := json.Marshal(LoginCheckResponse{
		Username:  username,
		AuthToken: authToken,
	})
	if err != nil {
		return err
	}

	return ctx.Send(loginCheckResp)
}

func insertCreds(db *pg.DB, username string, password string) error {
	_, err := db.Model(&model.UsersTable{
		username,
		password,
	}).Insert()
	return err
}

func validateCreds(db *pg.DB, username string, password string) error {
	isUserRegistered, err := db.Model(&model.UsersTable{}).Where("username = ? and password = ?", username, password).Count()
	if err != nil {
		return err
	}
	if isUserRegistered != 1 {
		return errors.New("user not registered")
	}
	return nil
}

func sendSuccessResponse(err error, ctx *fiber.Ctx, successStatusCode int, failureStatusCode int) error {
	if err != nil {
		return err
	} else {
		successResponse, err := json.Marshal(LoginRegisterResponse{
			IsRequestSuccessful: true,
		})
		if err != nil {
			return err
		}
		return ctx.Status(successStatusCode).Send(successResponse)
	}
}

func storeAuthCookieAndUpdateDatabase(ctx *fiber.Ctx, db *pg.DB) error {
	loginReq := LoginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return err
	}

	store, err := store.Sessions.Get(ctx) // get/create new session
	if err != nil {
		return err
	}
	err = store.Save()
	if err != nil {
		return err
	}

	_, err = db.Model(&model.CredsTable{
		UUID:     store.ID(),
		Username: loginReq.Username,
	}).Insert()
	return err
}
