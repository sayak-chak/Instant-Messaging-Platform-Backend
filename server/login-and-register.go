package server

import (
	// "database/sql"
	"database/sql"
	"errors"
	"fmt"

	// "github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	_ "github.com/lib/pq"
)
//TODO: refactor
var user = "USER"
var dbName="DB_NAME"
var postgresConfig = "user="+user+" dbname="+dbName+" sslmode=disable" 
var tableName = "TEST" 

type response struct {
	IsRequestSuccessful bool `json:"isRequestSuccessful"`
}

func register(ctx *websocket.Conn) {
	db, err := sql.Open("postgres", postgresConfig)
	if err != nil {
		fmt.Println("Register error => ", err)
	}
	err = insertCreds(db, ctx.Query("username"), ctx.Query("password"))
	if err != nil {
		fmt.Println("Register error => ", err)
		ctx.WriteJSON(response{
			IsRequestSuccessful: false,
		})
	} else {
		ctx.WriteJSON(response{
			IsRequestSuccessful: true,
		})
	}
	ctx.Close()
	defer db.Close()
}

func login(ctx *websocket.Conn) {
	db, err := sql.Open("postgres", postgresConfig) 
	if err != nil {
		fmt.Println("Register error => ", err)
	}
	err = validateCreds(db, ctx.Query("username"), ctx.Query("password"))
	if err != nil {
		fmt.Println("Register error => ", err)
		ctx.WriteJSON(response{
			IsRequestSuccessful: false,
		})
	} else {
		ctx.WriteJSON(response{
			IsRequestSuccessful: true,
		})
	}
	ctx.Close()
	defer db.Close()
}

func insertCreds(db *sql.DB, username string, password string) error {
	sqlres, err := db.Exec("create table if not exists "+tableName+"(username varchar(20), password varchar(20), PRIMARY KEY (username))")
	fmt.Println(sqlres)
	if err != nil {
		return err
	}

	sqlres, err = db.Exec("insert into "+tableName+" values ('" + username + "','" + password + "')")
	return err
}

func validateCreds(db *sql.DB, username string, password string) error {
	sqlres, err := db.Exec("select password from TEST where username='" + username + "' and password='" + password + "'")
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
