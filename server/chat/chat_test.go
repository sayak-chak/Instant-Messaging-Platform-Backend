package chat_test

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database/model"
	"instant-messaging-platform-backend/mocks"
	"instant-messaging-platform-backend/server/chat"
	"testing"
)

func Test_When_A_User_Sends_A_Message_To_Another_User_Then_It_Should_Be_Updated_In_Database(t *testing.T) {
	const SENDER = "SENDER"
	const RECIPIENT = "RECIPIENT"
	const MESSAGE = "MESSAGE"
	var numberOfTimesMessageIsSupposedToBeSentToReciepent = 3
	mocks.SetupDataBase()
	defer mocks.TearDown()
	mocks.RegisterUser(SENDER)
	mocks.RegisterUser(RECIPIENT)
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	for i := 0; i < numberOfTimesMessageIsSupposedToBeSentToReciepent; i++ {
		chat.SendDirectMessageTo(RECIPIENT,
			&chat.ChatResponse{
				Type:    "Personal",
				Sender:  SENDER,
				Message: MESSAGE,
			})
	}

	numberOfTimesMessageWasSentToReciepent, err := db.Model(&model.ChatTable{}).Where("sender = ? and message = ? and room = ?", SENDER, MESSAGE, "SENDER_RECIPIENT").Count()
	if err != nil {
		fmt.Print("ERROR ", err)
		t.Fail()
	}
	if numberOfTimesMessageWasSentToReciepent != numberOfTimesMessageIsSupposedToBeSentToReciepent {
		t.Fail()
	}
}
