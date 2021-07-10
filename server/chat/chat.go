package chat

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
	"instant-messaging-platform-backend/client"
	"instant-messaging-platform-backend/config"
	"instant-messaging-platform-backend/database/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gopkg.in/square/go-jose.v2/json"
)

var broadcastChannelForGroupChat = make(chan chatResponse)

type chatRequest struct {
	Target  string `json:"target"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type chatResponse struct {
	Type    string `json:"type"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type ChatHistoryResponses struct {
	Responses []chatHistoryResponse `json:"senderAndMessages"`
}

type chatHistoryResponse struct {
	Sender      string `json:"sender"`
	ChatMessage string `json:"message"`
}

func Chat(ctx *websocket.Conn) {
	//TODO: implement authentication
	if _, isClientAlreadyPresent := client.ClientsMap[ctx.Query("sender")]; isClientAlreadyPresent {
		return
	}
	client.AddClient(ctx.Query("sender"), ctx)
	defer ctx.Close()

	go broadcastMsgInGroup()

	for {
		chatRequestDetails := chatRequest{}
		var chatResponseDetails chatResponse

		err := ctx.ReadJSON(&chatRequestDetails)
		if err != nil { //TODO : add idle timer
			//TODO: figure out if closing the client is mandatory
			//err = client.ClientsMap[chatRequestDetails.Sender].Close()
			//if err != nil {
			//	print("Can't close client..", err)
			//}
			client.RemoveClient(chatRequestDetails.Sender)
			break
		}
		typeOfResponse := "Group"
		if chatRequestDetails.Target != typeOfResponse {
			typeOfResponse = "Personal"
		}
		chatResponseDetails = chatResponse{
			Type:    typeOfResponse,
			Sender:  chatRequestDetails.Sender,
			Message: chatRequestDetails.Message,
		}
		if typeOfResponse == "Group" {
			broadcastChannelForGroupChat <- chatResponseDetails
		} else {
			sendDirectMessageTo(chatRequestDetails.Target, chatResponseDetails)
		}
	}
}

func sendDirectMessageTo(target string, chatResponseDetails chatResponse) {
	targetCtx, isTargetClientPresent := client.ClientsMap[target]
	if isTargetClientPresent {
		err := targetCtx.WriteJSON(chatResponseDetails)
		if err != nil {
			fmt.Println(err)
		}
		ownCxt, isPresent := client.ClientsMap[chatResponseDetails.Sender]
		if !isPresent {
			storeMessageInDataBase(chatResponseDetails.Sender, chatResponseDetails.Message, target)
			return
		}
		ownCxt.WriteJSON(chatResponseDetails)
		storeMessageInDataBase(chatResponseDetails.Sender, chatResponseDetails.Message, target)
		return

	}
	ownCxt, isPresent := client.ClientsMap[chatResponseDetails.Sender]
	if !isPresent {
		storeMessageInDataBase(chatResponseDetails.Sender, chatResponseDetails.Message, target)
		return
	}

	ownCxt.WriteJSON(chatResponseDetails)
	storeMessageInDataBase(chatResponseDetails.Sender, chatResponseDetails.Message, target)
}

func broadcastMsgInGroup() {
	for {
		msg := <-broadcastChannelForGroupChat

		for specificClient := range client.ClientsMap {
			err := client.ClientsMap[specificClient].WriteJSON(msg)
			if err != nil { // close that client & remove it
				client.ClientsMap[specificClient].Close()
				client.RemoveClient(specificClient)
			}
		}
	}
}

func storeMessageInDataBase(sender string, message string, target string) {
	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()
	room := getRoomName(sender, target)
	_, err := db.Model(&model.ChatTable{
		sender,
		room,
		message,
	}).Insert()
	if err != nil {
		fmt.Println(err)
	}
}

func GetChatHistory(ctx *fiber.Ctx) error {
	//TODO: authneticate
	username := ctx.Query("username")
	sender := ctx.Query("sender")

	db := pg.Connect(&pg.Options{
		User:     config.User,
		Database: config.DbName,
	})
	defer db.Close()

	chatSenderAndMessageList := make([]chatHistoryResponse, 0)

	var senderUsername []string
	var messageToBeSent []string

	room := getRoomName(username, sender)
	count, err := db.Model(&model.ChatTable{}).ColumnExpr("array_agg(sender), array_agg(message)").Where("room = ?", room).SelectAndCount(pg.Array(&senderUsername), pg.Array(&messageToBeSent))

	if err != nil {
		return err
	}

	for idx := 0; idx < count; idx++ {
		chatSenderAndMessageList = append(chatSenderAndMessageList, chatHistoryResponse{
			Sender:      senderUsername[idx],
			ChatMessage: messageToBeSent[idx],
		})
	}

	chatHistoryResponse, err := json.Marshal(ChatHistoryResponses{
		Responses: chatSenderAndMessageList,
	})
	if err != nil {
		return err
	}

	return ctx.Status(200).Send(chatHistoryResponse)
}

func getRoomName(userOne string, userTwo string) string {
	if userOne > userTwo {
		return userOne + "_" + userTwo
	}
	return userTwo + "_" + userOne
}
