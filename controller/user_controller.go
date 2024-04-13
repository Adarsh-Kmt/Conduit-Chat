package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adarsh-Kmt/chatapp/service"
	"github.com/Adarsh-Kmt/chatapp/types"
	"github.com/Adarsh-Kmt/chatapp/util"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
)

type UserController struct {
	UserService service.UserService
}

func NewUserControllerInstance(UserServiceInstance service.UserService) *UserController {

	return &UserController{
		UserService: UserServiceInstance,
	}
}
func (userController *UserController) Run(router *mux.Router) *mux.Router {

	router.HandleFunc("/register", util.MakeHttpHandlerFunc(userController.HandleRegisterUser))
	router.HandleFunc("/chatRequest", util.MakeJWTAuthHttpHandlerFunc(util.MakeHttpHandlerFunc(userController.HandleSendChatRequest)))
	router.HandleFunc("/login", util.MakeHttpHandlerFunc(userController.HandleLoginUser))
	router.HandleFunc("/incomingChatRequest", util.MakeJWTAuthHttpHandlerFunc(util.MakeHttpHandlerFunc(userController.GetIncomingChatRequestList)))
	router.HandleFunc("/outgoingChatRequest", util.MakeJWTAuthHttpHandlerFunc(util.MakeHttpHandlerFunc(userController.GetOutgoingChatRequestList)))

	return router
	//http.ListenAndServe(s.ListAddr, userrouter)
}

func (userController *UserController) HandleRegisterUser(w http.ResponseWriter, r *http.Request) *types.APIError {

	request := new(types.RegisterUserRequest)
	//var request bson.M
	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {
		return &types.APIError{Error: "error in parsing POST request body.", ErrorStatus: 500}
	}

	ApiError := userController.UserService.RegisterUser(request)

	if ApiError != nil {
		return ApiError
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"SuccessMessage": "you have successfully registered."})

	return nil

}

func (userController *UserController) HandleLoginUser(w http.ResponseWriter, r *http.Request) *types.APIError {

	request := new(types.LoginUserRequest)

	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {

		return &types.APIError{Error: "error in parsing POST request body.", ErrorStatus: 500}
	}

	JwtToken, ApiError := userController.UserService.LoginUser(request)

	if ApiError != nil {

		return ApiError
	}

	util.WriteJSON(w, 200, map[string]string{"SuccessMessage": JwtToken})
	return nil
}
func (userController *UserController) HandleSendChatRequest(w http.ResponseWriter, r *http.Request) *types.APIError {

	fmt.Println("entered the controller.")
	request := new(types.ChatRequest)

	JwtToken := r.Header.Get("Auth")

	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {

		return &types.APIError{Error: "error in parsing POST request body.", ErrorStatus: 500}
	}

	fmt.Println(request.SenderId.Hex())

	ValidatedUserObjectIdString, err := util.GetUserObjectIdFromJWT(JwtToken)

	fmt.Println("user object id is:" + ValidatedUserObjectIdString)

	if err != nil {
		return &types.APIError{Error: err.Error(), ErrorStatus: 500}
	}
	ValidatedUserObjectId, er := primitive.ObjectIDFromHex(ValidatedUserObjectIdString)

	if er != nil {

		return &types.APIError{Error: er.Error(), ErrorStatus: 500}
	}

	ChatRequestObjectId, ApiError := userController.UserService.SendChatRequest(ValidatedUserObjectId, request)

	if ApiError != nil {
		return ApiError
	}

	util.WriteJSON(w, http.StatusOK, map[string]string{"SuccessMessage": ChatRequestObjectId})

	return nil
}

func (userController *UserController) GetIncomingChatRequestList(w http.ResponseWriter, r *http.Request) *types.APIError {

	JwtToken := r.Header.Get("Auth")

	ValidatedUserObjectIdString, err := util.GetUserObjectIdFromJWT(JwtToken)

	if err != nil {
		return &types.APIError{Error: err.Error(), ErrorStatus: 500}
	}

	ValidatedUserObjectId, err := primitive.ObjectIDFromHex(ValidatedUserObjectIdString)

	if err != nil {

		return &types.APIError{Error: "internal error 500", ErrorStatus: 500}
	}

	IncomingChatRequestList, ApiError := userController.UserService.GetIncomingChatRequestList(ValidatedUserObjectId)

	if ApiError != nil {

		return ApiError
	}

	util.WriteJSON(w, 200, map[string][]types.ChatRequest{"SuccessMessage": IncomingChatRequestList})

	return nil

}

func (userController *UserController) GetOutgoingChatRequestList(w http.ResponseWriter, r *http.Request) *types.APIError {

	JwtToken := r.Header.Get("Auth")

	ValidatedUserObjectIdString, err := util.GetUserObjectIdFromJWT(JwtToken)

	if err != nil {
		return &types.APIError{Error: err.Error(), ErrorStatus: 500}
	}

	ValidatedUserObjectId, err := primitive.ObjectIDFromHex(ValidatedUserObjectIdString)

	OutgoingChatRequestList, ApiError := userController.UserService.GetOutgoingChatRequestList(ValidatedUserObjectId)

	if ApiError != nil {

		return ApiError
	}

	util.WriteJSON(w, 200, map[string][]types.ChatRequest{"SuccessMessage": OutgoingChatRequestList})

	return nil

}
