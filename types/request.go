package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AcceptChatRequest struct {
	ChatRequestId primitive.ObjectID `json:"chatRequestId"`
}

type MessageRequest struct {
	ReceiverId primitive.ObjectID `json:"receiverId"`
	Message    string             `json:"message"`
}
