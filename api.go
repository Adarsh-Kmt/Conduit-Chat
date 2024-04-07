package main

import (
	"encoding/json"
	"fmt"
	"log"

	// "fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type APIServer struct {
	ListAddr   string
	Storage    store
	activeConn map[string]*websocket.Conn
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
	router.HandleFunc("/register", MakeHttpHandlerFunc(s.HandleRegisterUser))
	router.HandleFunc("/chatRequest", MakeJWTAuthHttpHandlerFunc(MakeHttpHandlerFunc(s.HandlerSendChatRequest)))
	router.HandleFunc("/message", MakeJWTAuthHttpHandlerFunc(MakeHttpHandlerFunc(s.HandleSendMessage)))
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

	UserObjectId, err := getUserObjectIdFromJWT(JWTToken)

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
		var message MessageRequest

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

func (s *APIServer) HandleRegisterUser(w http.ResponseWriter, r *http.Request) *APIError {

	request := new(RegisterUserRequest)
	//var request bson.M
	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {

		return &APIError{Error: "error in parsing POST request body.", ErrorStatus: 500}
	}

	newUser := CreateNewUser(request)

	JWTToken, ApiError := s.Storage.RegisterUser(newUser)

	if ApiError != nil {
		return ApiError
	}

	WriteJSON(w, http.StatusOK, JWTToken)

	return nil

}

func (s *APIServer) HandlerSendChatRequest(w http.ResponseWriter, r *http.Request) *APIError {

	request := new(ChatRequest)

	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {

		return &APIError{Error: "error in parsing POST request body.", ErrorStatus: 500}
	}

	ChatRequestObjectId, ApiError := s.Storage.SendChatRequest(request)

	if ApiError != nil {

		return ApiError
	}

	WriteJSON(w, http.StatusOK, ChatRequestObjectId)

	return nil
}

type ApiFunc func(http.ResponseWriter, *http.Request) *APIError

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
		_, err := ValidateJWTToken(tokenString)

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
