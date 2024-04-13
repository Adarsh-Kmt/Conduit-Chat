package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                      primitive.ObjectID `json:"objectId" bson:"_id, omitempty"`
	Name                    string             `json:"name" bson:"Name"`
	Password                string             `json:"password" bson:"Password"`
	ChatSummaryList         []ChatSummary      `json:"chatIdList" bson:"ChatIdList"`
	UserConnList            map[string]string  `json:"userConnList" bson:"UserConnList"`
	OutgoingChatRequestList []ChatRequest      `json:"outgoingChatRequestList" bson:"OutgoingChatRequestList"`
	IncomingChatRequestList []ChatRequest      `json:"incomingChatRequestList" bson:"IncomingChatRequestList"`
}

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Chat struct {
	ChatId      primitive.ObjectID `json:"chatId" bson:"_id,omitempty"`
	User1       primitive.ObjectID `json:"user1ObjectId" bson:"User1ObjectId"`
	User2       primitive.ObjectID `json:"user2ObjectId" bson:"User2ObjectId"`
	MessageList []Message          `json:"messageList" bson:"MessageList"`
}

type ChatSummary struct {
	ChatSummaryId   primitive.ObjectID `json:"chatSummaryId" bson:"_id,omitempty"`
	User            primitive.ObjectID `json:"userObjectId" bson:"UserObjectId"`
	ChatReferenceId primitive.ObjectID `json:"chatReferenceId" bson:"ChatReferenceId"`
}
type Message struct {
	MessageId  primitive.ObjectID `json:"messageId" bson:"_id,omitempty"`
	SenderId   primitive.ObjectID `json:"senderId" bson:"SenderObjectId"`
	ReceiverId primitive.ObjectID `json:"receiverId" bson:"ReceiverObjectId"`
	Message    string             `json:"message" bson:"Message"`
	SentAtTime time.Time          `json:"sentAtTime" bson:"SentAtTime"`
}

type ChatRequest struct {
	ChatRequestId   primitive.ObjectID `json:"chatRequestId" bson:"_id,omitempty"`
	SenderId        primitive.ObjectID `json:"senderId" bson:"SenderObjectId"`
	ReceiverId      primitive.ObjectID `json:"receiverId" bson:"ReceiverObjectId"`
	SentOnDate      time.Time          `json:"sentOnDate" bson:"SentOnDate"`
	RequestAccepted bool               `json:"requestAccepted" bson:"RequestAccepted"`
}

type MessageRequest struct {
	ReceiverId primitive.ObjectID `json:"receiverId"`
	Message    string             `json:"message"`
}

type APIError struct {
	Error       string
	ErrorStatus int
}

type OutgoingChatRequestListResponse struct {
	OutgoingChatRequestList []ChatRequest `json:"outgoingChatRequestList" bson:"OutgoingChatRequestList"`
}

type IncomingChatRequestListResponse struct {
	IncomingChatRequestList []ChatRequest `json:"incomingChatRequestList" bson:"IncomingChatRequestList"`
}
