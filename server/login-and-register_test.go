package server_test

import (
	"bytes"
	"database/sql"
	"fmt"

	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/mocks"
	"instant-messaging-platform-backend/server"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	utils "github.com/gofiber/utils"
)

func Test_When_Guest_Tries_To_Create_New_User_With_Unique_Username_Then_Server_Should_Return_Status_201(t *testing.T) {
	defer tearDown()
	app := server.New()
	app.Listen("3000")

	username := "TEST_USERNAME"
	password := "TEST_PASSWORD"
	jsonStr, err := mocks.GetLoginJson(username, password)
	if err != nil {
		t.Error(err)
	}
	request := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, `{"isRequestSuccessful":true}`, string(responseBody), "Response Body")
	utils.AssertEqual(t, 201, resp.StatusCode, "Status code")
}

func Test_When_Guest_Tries_To_Create_New_User_With_Non_Unique_Username_Then_Server_Should_Return_Status_409(t *testing.T) {
	defer tearDown()
	app := server.New()
	app.Listen("3000")
	nonUniqueUsername := "TEST_USERNAME"
	nonUniquePassword := "TEST_PASSWORD"
	jsonStr, err := mocks.GetLoginJson(nonUniqueUsername, nonUniquePassword)
	if err != nil {
		t.Error(err)
	}
	createUser(nonUniqueUsername, nonUniquePassword, jsonStr, app, t)

	request := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, `{"isRequestSuccessful":false}`, string(responseBody), "Response Body")
	utils.AssertEqual(t, 409, resp.StatusCode, "Status code")
}

func Test_When_Guest_Tries_To_Login_With_Valid_Credentials_Then_Server_Should_Return_Status_200(t *testing.T) {
	defer tearDown()
	app := server.New()
	app.Listen("3000")
	username := "TEST_USERNAME"
	password := "TEST_PASSWORD"
	jsonStr, err := mocks.GetLoginJson(username, password)
	if err != nil {
		t.Error(err)
	}
	createUser(username, password, jsonStr, app, t)

	request := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, `{"isRequestSuccessful":true}`, string(responseBody), "Response Body")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}
func Test_When_Guest_Tries_To_Login_With_Invalid_Credentials_Then_Server_Should_Return_Status_400(t *testing.T) {
	defer tearDown()
	app := server.New()
	app.Listen("3000")
	username := "TEST_USERNAME"
	password := "TEST_PASSWORD"
	invalidUserName := "INVALID_USERNAME"
	jsonStr, err := mocks.GetLoginJson(username, password)
	if err != nil {
		t.Error(err)
	}
	createUser(username, password, jsonStr, app, t)
	jsonStr, err = mocks.GetLoginJson(invalidUserName, password)
	if err != nil {
		t.Error(err)
	}

	request := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	
	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, `{"isRequestSuccessful":false}`, string(responseBody), "Response Body")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func createUser(username string, password string, jsonStr []byte, app *fiber.App, t *testing.T){
	request := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")

	_, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	jsonStr, err = mocks.GetLoginJson(username, password)
	if err != nil {
		t.Error(err)
	}
}

func tearDown() {
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec("drop table " + config.CredsTable)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("drop table " + config.UsersTable)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("drop table " + config.ChatTable)
	if err != nil {
		fmt.Println(err)
	}
}