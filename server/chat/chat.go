package chat

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"instant-messaging-platform-backend/client"
	"instant-messaging-platform-backend/config"

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
			client.ClientsMap[chatRequestDetails.Sender].Close()
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

func GetChatHistory(ctx *fiber.Ctx) error {
	//TODO: authneticate
	username := ctx.Query("username")
	sender := ctx.Query("sender")

	db, err := sql.Open("postgres", config.PostgresConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var message string
	chatSenderAndMessageList := make([]chatHistoryResponse, 0)

	rowPtrs, err := db.Query("select sender,message from " + config.ChatTable + " where room='" + getRoomName(username, sender) + "'")

	if err != nil {
		return err
	}

	for rowPtrs.Next() {
		if err := rowPtrs.Scan(&sender, &message); err != nil {
			return err
		}
		chatSenderAndMessageList = append(chatSenderAndMessageList, chatHistoryResponse{
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
