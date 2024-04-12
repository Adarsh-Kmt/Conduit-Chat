package data

import (
	"time"

	"github.com/Adarsh-Kmt/chatapp/types"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewUser(request *types.RegisterUserRequest) *types.User {

	return &types.User{
		Name:                    request.Name,
		Password:                request.Password,
		ChatSummaryList:         []types.ChatSummary{},
		UserConnList:            map[string]string{},
		OutgoingChatRequestList: []types.ChatRequest{},
		IncomingChatRequestList: []types.ChatRequest{},
	}
}

func CreateUserBSON(NewUser *types.User) bson.D {

	return bson.D{
		{Key: "Name", Value: NewUser.Name},
		{Key: "Password", Value: NewUser.Password},
		{Key: "ChatIdList", Value: NewUser.ChatSummaryList},
		{Key: "UserConnList", Value: NewUser.UserConnList},
		{Key: "OutgoingChatRequestList", Value: NewUser.OutgoingChatRequestList},
		{Key: "IncomingChatRequestList", Value: NewUser.IncomingChatRequestList},
	}
}

func CreateChatRequestBSON(request *types.ChatRequest) (string, bson.D) {

	chatRequestId := primitive.NewObjectID()
	return chatRequestId.Hex(), bson.D{
		{Key: "_id", Value: chatRequestId},
		{Key: "SenderObjectId", Value: request.SenderId},
		{Key: "ReceiverObjectId", Value: request.ReceiverId},
		{Key: "SentOnDate", Value: time.Now()},
		{Key: "RequestAccepted", Value: false},
	}
}
