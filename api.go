package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Adarsh-Kmt/chatapp/types"
	"github.com/Adarsh-Kmt/chatapp/util"

	"github.com/Adarsh-Kmt/chatapp/controller"
	// "fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

/*
flow:

if A wants to send a message to B, message Request must include chatId and message string.

2 tabs in application, one for chats, and other for chat requests.

if A sends chat request to B, and B accepts, Both must be returned the chatObjectId corresponding to the chat.

this is simple for the user that accepts the request. a response to the accept message is the chatId.

for the person that sends the chat request, he will need to go to the chat request tab to see if the request has been accepted,
if chat request has been accepted, clicking on the chat request gets him the chatId,
otherwise response should be sent asking him to wait for receiver to accept request.
*/
type APIServer struct {
	UserController controller.UserController
	ListAddr       string
	Storage        store
	activeConn     map[string]*websocket.Conn
}

func NewAPIServer(listAddr string, storage store) *APIServer {

	return &APIServer{

		ListAddr:   listAddr,
		Storage:    storage,
		activeConn: make(map[string]*websocket.Conn),
	}
}

func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/register", MakeHttpHandlerFunc(s.UserController.HandleRegisterUser))
	//router.HandleFunc("/chatRequest", MakeJWTAuthHttpHandlerFunc(MakeHttpHandlerFunc(sUser.HandleSendChatRequest)))
	//router.HandleFunc("/message", MakeJWTAuthHttpHandlerFunc(MakeHttpHandlerFunc(s.HandleSendMessage)))
	http.ListenAndServe(s.ListAddr, router)
}

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// error handling does not work for websockets.
// cannot use WRITEJSON method, as after upgrade, switching to websocket protocol,
// http is no longer used. so we cannot use http.ResponseWriter to write a response back to the user. need to fix that issue.
func (s *APIServer) HandleSendMessage(w http.ResponseWriter, r *http.Request) *APIError {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {

		return &APIError{Error: "error while establishing websocket connection.", ErrorStatus: 500}
	}

	JWTToken := r.Header.Get("Auth")

	UserObjectId, err := util.GetUserObjectIdFromJWT(JWTToken)

	if err != nil {
		log.Println(err.Error())
		fmt.Println("error while getting userobejctId from jwt.")

		// this return statement is an isse, results in the WRITEJSON method trying to use http.ResponseWriter
		return &APIError{Error: err.Error(), ErrorStatus: 500}
	}

	s.activeConn[UserObjectId] = conn

	for {

		_, messageByteArray, err := conn.ReadMessage()

		if err != nil {
			fmt.Println(err.Error())

			// this return statement is an isse, results in the WRITEJSON method trying to use http.ResponseWriter
			return &APIError{Error: err.Error(), ErrorStatus: 500}

		}
		var message types.MessageRequest

		if err := json.Unmarshal([]byte(messageByteArray), &message); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		fmt.Println("message ", message.Message, " sent to receiver:", message.ReceiverId)

		ReceiverConn, exists := s.activeConn[message.ReceiverId.Hex()]

		if !exists {
			conn.WriteMessage(websocket.TextMessage, []byte(message.ReceiverId.Hex()+" is not online."))
		} else {
			ReceiverConn.WriteMessage(websocket.TextMessage, []byte(message.Message))
		}
		//conn.WriteMessage(websocket.TextMessage, []byte("message "+message.Message+" sent to receiver:"+message.ReceiverId.Hex()))
	}

}

type ApiFunc func(http.ResponseWriter, *http.Request) *types.APIError

type APIError struct {
	Error       string
	ErrorStatus int
}

func WriteJSON(w http.ResponseWriter, status int, body any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

func MakeJWTAuthHttpHandlerFunc(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Auth")

		// if substr := tokenString[:7]; len(tokenString) == 0 || substr != "Bearer " {

		// 	WriteJSON(w, http.StatusForbidden, map[string]string{"Error": "incorrect auth token."})
		// 	return

		// }
		log.Println(tokenString)
		_, err := util.ValidateJWTToken(tokenString)

		if err != nil {

			WriteJSON(w, http.StatusForbidden, map[string]string{"Error": err.Error()})
			return
		}

		f(w, r)

	}
}

func MakeHttpHandlerFunc(f ApiFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if APIError := f(w, r); APIError != nil {

			WriteJSON(w, APIError.ErrorStatus, map[string]string{"error": APIError.Error})
		}
	}
}
