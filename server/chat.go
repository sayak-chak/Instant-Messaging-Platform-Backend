package server

import (
	"database/sql"
	"fmt"
	"instant-messaging-platform-backend/client"
	"instant-messaging-platform-backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gopkg.in/square/go-jose.v2/json"
)

var broadcastChannelForGroupChat = make(chan ChatResponse)

type ChatRequest struct {
	Target  string `json:"target"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type ChatResponse struct {
	Type    string `json:"type"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type ChatHistoryResponses struct {
	Responses []ChatHistoryResponse `json:"senderAndMessages"`
}

type ChatHistoryResponse struct {
	Sender      string `json:"sender"`
	ChatMessage string `json:"message"`
}

func chat(ctx *websocket.Conn) {
	//TODO: implement authentication
	if _, isClientAlreadyPresent := client.ClientsMap[ctx.Query("sender")]; isClientAlreadyPresent {
		return
	}
	client.AddClient(ctx.Query("sender"), ctx)
	defer ctx.Close()

	go broadcastMsgInGroup()

	for {
		chatRequest := ChatRequest{}
		var chatResponse ChatResponse

		err := ctx.ReadJSON(&chatRequest)
		if err != nil { //TODO : add idle timer
			client.ClientsMap[chatRequest.Sender].Close()
			client.RemoveClient(chatRequest.Sender)
			break
		}
		typeOfResponse := "Group"
		if chatRequest.Target != typeOfResponse {
			typeOfResponse = "Personal"
		}
		chatResponse = ChatResponse{
			Type:    typeOfResponse,
			Sender:  chatRequest.Sender,
			Message: chatRequest.Message,
		}
		if typeOfResponse == "Group" {
			broadcastChannelForGroupChat <- chatResponse
		} else {
			sendDirectMessageTo(chatRequest.Target, chatResponse)
		}
	}
}

func sendDirectMessageTo(target string, chatResponse ChatResponse) {
	targetCtx, isTargetClientPresent := client.ClientsMap[target]
	if isTargetClientPresent {
		err := targetCtx.WriteJSON(chatResponse)
		if err != nil {
			fmt.Println(err)
		}
		ownCxt, isPresent := client.ClientsMap[chatResponse.Sender]
		if !isPresent {
			storeMessageInDataBase(chatResponse.Sender, chatResponse.Message, target)
			return
		}
		ownCxt.WriteJSON(chatResponse)
		storeMessageInDataBase(chatResponse.Sender, chatResponse.Message, target)
		return

	}
	ownCxt, isPresent := client.ClientsMap[chatResponse.Sender]
	if !isPresent {
		storeMessageInDataBase(chatResponse.Sender, chatResponse.Message, target)
		return
	}

	ownCxt.WriteJSON(chatResponse)
	storeMessageInDataBase(chatResponse.Sender, chatResponse.Message, target)
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
	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	room := getRoomName(sender, target)
	_, err = db.Query("insert into " + config.ChatTable + " values ('" + sender + "', '" + room + "', '" + message + "')")
	if err != nil {
		fmt.Println(err)
	}
}

func getChatHistory(ctx *fiber.Ctx) error {
	//TODO: authneticate
	username := ctx.Query("username")
	sender := ctx.Query("sender")

	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var message string
	chatSenderAndMessageList := make([]ChatHistoryResponse, 0)

	rowPtrs, err := db.Query("select sender,message from " + config.ChatTable + " where room='" + getRoomName(username, sender) + "'")

	if err != nil {
		return err
	}

	for rowPtrs.Next() {
		if err := rowPtrs.Scan(&sender, &message); err != nil {
			return err
		}
		chatSenderAndMessageList = append(chatSenderAndMessageList, ChatHistoryResponse{
			Sender:      sender,
			ChatMessage: message,
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
