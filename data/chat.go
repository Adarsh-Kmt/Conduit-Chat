package data

import (
	"github.com/Adarsh-Kmt/chatapp/types"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewChatBSON(request types.ChatRequest) (primitive.ObjectID, bson.D) {

	chatObjectId := primitive.NewObjectID()

	return chatObjectId, bson.D{
		{Key: "ChatId", Value: chatObjectId},
		{Key: "User1", Value: request.SenderId},
		{Key: "User2", Value: request.ReceiverId},
		{Key: "MessageList", Value: []types.Message{}},
	}

}
