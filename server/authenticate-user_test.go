package server_test

import (
	"fmt"
	"instant-messaging-platform-backend/database"
	"instant-messaging-platform-backend/mocks"
	"instant-messaging-platform-backend/server"
	"testing"
)

func Test_When_A_Valid_Auth_Token_Is_Provided_Then_Is_Auth_Valid_Should_Return_True(t *testing.T) {
	err := database.SetupDataBase()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer tearDown()
	validUUID := "9c7eb310-bd8f-11eb-8529-0242ac130003"
	username := "TEST_USERNAME"
	err = mocks.FillAuthTokenInDB(validUUID, username)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	isAuthValid, err := server.IsAuthValid(validUUID)

	if isAuthValid != true {
		t.Fail()
	}
}

func Test_When_An_InValid_Auth_Token_Is_Provided_Then_Is_Auth_Valid_Should_Return_False(t *testing.T) {
	err := database.SetupDataBase()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer tearDown()
	validUUID := "9c7eb310-bd8f-11eb-8529-0242ac130003"
	invalidUUID := "8c7eb310-bd8f-11eb-8529-0242ac130003"
	username := "TEST_USERNAME"
	err = mocks.FillAuthTokenInDB(validUUID, username)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	isAuthValid, err := server.IsAuthValid(invalidUUID)

	if isAuthValid != false {
		t.Fail()
	}
}
