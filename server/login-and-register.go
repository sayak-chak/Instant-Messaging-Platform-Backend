package server

import (
	// "database/sql"
	"database/sql"
	"errors"
	"instant-messaging-platform-backend/config"

	"encoding/json"

	"github.com/gofiber/fiber/v2"

	// "github.com/gofiber/websocket/v2"
	_ "github.com/lib/pq"
)

type LoginCheckResponse struct {
	Username string `json:"username"`
}
type LoginRegisterResponse struct {
	IsRequestSuccessful bool `json:"isRequestSuccessful"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(ctx *fiber.Ctx) error {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	loginReq := LoginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return err
	}

	err = sendResponse(insertCreds(db, loginReq.Username, loginReq.Password), ctx, 201, 409)
	if err != nil {
		return nil //TODO: unless err == nil, it wasn't giving custom status code, maybe refactor
	}

	return storeAuthCookieAndUpdateDatabase(ctx, db)
}

func Login(ctx *fiber.Ctx) error {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	loginReq := LoginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return err
	}
	
	err = sendResponse(validateCreds(db, loginReq.Username, loginReq.Password), ctx, 200, 400)
	if err != nil {
		return nil //TODO: unless err == nil, it wasn't giving custom status code, maybe refactor
	}

	return storeAuthCookieAndUpdateDatabase(ctx, db)
}

func LoginCheck(ctx *fiber.Ctx) error {
	var uuid string
	var username string

	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	rowPtr := db.QueryRow("select * from " + config.CredsTable + " where uuid='" + string(ctx.Request().Header.Cookie("Auth-Token")) + "'")
	if err := rowPtr.Scan(&uuid, &username); err != nil {
		return err
	}

	loginCheckResp, err := json.Marshal(LoginCheckResponse{
		Username: username,
	})
	if err != nil {
		return err
	}

	return ctx.Send(loginCheckResp)
}

func insertCreds(db *sql.DB, username string, password string) error {
	_, err := db.Exec("create table if not exists " + config.UsersTable + "(username varchar(20), password varchar(20), PRIMARY KEY (username))")
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into " + config.UsersTable + " values ('" + username + "','" + password + "')")
	return err
}

func validateCreds(db *sql.DB, username string, password string) error {
	sqlres, err := db.Exec("select password from " + config.UsersTable + " where username='" + username + "' and password='" + password + "'")
	if err != nil {
		return err
	}
	isUserRegistered, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}
	if isUserRegistered != 1 {
		return errors.New("User not registered")
	}
	return nil
}

func sendResponse(err error, ctx *fiber.Ctx, successStatusCode int, failureStatusCode int) error {
	if err != nil {
		failureResponse, error := json.Marshal(LoginRegisterResponse{
			IsRequestSuccessful: false,
		})
		if error != nil {
			return error
		}
		error = ctx.Send(failureResponse)
		if error != nil {
			return error
		}
		error = ctx.Status(failureStatusCode).Send(failureResponse)
		if error != nil {
			return error
		}
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

func storeAuthCookieAndUpdateDatabase(ctx *fiber.Ctx, db *sql.DB) error {
	loginReq := LoginRequest{}
	if err := ctx.BodyParser(&loginReq); err != nil {
		return err
	}

	store, err := Sessions.Get(ctx) // get/create new session
	if err != nil {
		return err
	}

	_, err = db.Exec("create table if not exists " + config.CredsTable + "(uuid varchar(36), username varchar(20), PRIMARY KEY (uuid))") //varchar 36 to store uuid
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into " + config.CredsTable + " values ('" + store.ID() + "','" + loginReq.Username + "')")
	return err
}
