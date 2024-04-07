package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                     primitive.ObjectID `json:"objectId" bson:"_id"`
	Name                   string             `json:"name" bson:"Name"`
	ChatIdList             map[string]string  `json:"chatIdList" bson:"ChatIdList"`
	UserConnList           map[string]string  `json:"userConnList" bson:"UserConnList"`
	PendingChatRequestList map[string]string  `json:"pendingChatRequestList" bson:"PendingChatRequestList"`
}

type RegisterUserRequest struct {
	Name string `json:"name"`
}

func CreateNewUser(request *RegisterUserRequest) *User {

	return &User{
		Name:                   request.Name,
		ChatIdList:             map[string]string{},
		UserConnList:           map[string]string{},
		PendingChatRequestList: map[string]string{},
	}
}

type ChatRequest struct {
	ChatRequestId primitive.ObjectID `json:"chatRequestId"`
	SenderId      primitive.ObjectID `json:"senderId"`
	ReceiverId    primitive.ObjectID `json:"receiverId"`
	SentOnDate    time.Time          `json:"sentOnDate"`
}

type MessageRequest struct {
	ReceiverId primitive.ObjectID `json:"receiverId"`
	Message    string             `json:"message"`
}
