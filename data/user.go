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

func CreateChatRequestBSONS(request *types.ChatRequest) (string, bson.D, string, bson.D) {

	senderChatRequestId := primitive.NewObjectID()
	receiverChatRequestId := primitive.NewObjectID()

	senderChatRequestBSON := bson.D{
		{Key: "_id", Value: senderChatRequestId},
		{Key: "ReferenceId", Value: receiverChatRequestId},
		{Key: "SenderObjectId", Value: request.SenderId},
		{Key: "ReceiverObjectId", Value: request.ReceiverId},
		{Key: "SentOnDate", Value: time.Now()},
		{Key: "RequestAccepted", Value: false},
	}

	receiverChatRequestBSON := bson.D{
		{Key: "_id", Value: receiverChatRequestId},
		{Key: "ReferenceId", Value: senderChatRequestId},
		{Key: "SenderObjectId", Value: request.SenderId},
		{Key: "ReceiverObjectId", Value: request.ReceiverId},
		{Key: "SentOnDate", Value: time.Now()},
		{Key: "RequestAccepted", Value: false},
	}

	return senderChatRequestId.Hex(), senderChatRequestBSON, receiverChatRequestId.Hex(), receiverChatRequestBSON
}
