package service

import (
	"encoding/json"

	"github.com/Adarsh-Kmt/chatapp/repository"
	"github.com/Adarsh-Kmt/chatapp/types"
	"github.com/gorilla/websocket"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	SendMessage(*websocket.Conn, string) *types.APIError
}

type MessageServiceImpl struct {
	UserRepository repository.UserRepository
	activeConn     map[string]*websocket.Conn
}

func MakeMessageServiceImplInstance(userRepositoryImpl repository.UserRepository) *MessageServiceImpl {

	return &MessageServiceImpl{
		UserRepository: userRepositoryImpl,
		activeConn:     make(map[string]*websocket.Conn),
	}
}

func (messageService *MessageServiceImpl) SendMessage(SenderConn *websocket.Conn, userObjectId string) *types.APIError {

	//distributorNode.HandleMessage(conn.getMessage)

	messageService.activeConn[userObjectId] = SenderConn

	for {

		_, message, err := SenderConn.ReadMessage()

		if err != nil {

			return &types.APIError{Error: err.Error(), ErrorStatus: 500}
		}

		var mr types.MessageRequest

		if err := json.Unmarshal(message, &mr); err != nil {

			return &types.APIError{Error: err.Error(), ErrorStatus: 500}
		}

		ReceiverConn, exists := messageService.activeConn[mr.ReceiverId.Hex()]

		if exists {

			ReceiverConn.WriteMessage(websocket.TextMessage, []byte(mr.Message))

		} else {

			SenderConn.WriteMessage(websocket.TextMessage, []byte("user is offline."))
		}

	}

}
