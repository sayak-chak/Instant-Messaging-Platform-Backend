package search_test

import (
	"instant-messaging-platform-backend/mocks"
	"instant-messaging-platform-backend/server"
	"instant-messaging-platform-backend/server/search"

	"net/http/httptest"
	"testing"

	utils "github.com/gofiber/utils"
)

func Test_When_Search_Is_Called_On_Valid_Username_Then_It_Should_Return_Status_200(t *testing.T) {
	defer mocks.TearDown()
	app := server.New()
	app.Listen("3000")
	search.AuthValidator = mocks.AuthValidator

	username := "TEST_USERNAME"
	err := mocks.RegisterUser(username)
	if err != nil {
		t.Error(err)
	}
	request := httptest.NewRequest("GET", "/search?username="+username, nil)

	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func Test_When_Search_Is_Called_On_Invalid_Username_Then_It_Should_Return_Status_400(t *testing.T) {
	defer mocks.TearDown()
	app := server.New()
	app.Listen("3000")
	search.AuthValidator = mocks.AuthValidator

	username := "TEST_USERNAME"
	err := mocks.RegisterUser(username)
	if err != nil {
		t.Error(err)
	}
	request := httptest.NewRequest("GET", "/search?username="+"INVALID_USERNAME", nil)

	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}
