package service

import (
	"fmt"

	"github.com/Adarsh-Kmt/chatapp/data"
	"github.com/Adarsh-Kmt/chatapp/repository"
	"github.com/Adarsh-Kmt/chatapp/types"
	"github.com/Adarsh-Kmt/chatapp/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	RegisterUser(*types.RegisterUserRequest) *types.APIError
	LoginUser(*types.LoginUserRequest) (string, *types.APIError)
	SendChatRequest(primitive.ObjectID, *types.ChatRequest) (string, *types.APIError)
	GetOutgoingChatRequestList(string) []types.ChatRequest
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewUserServiceImplInstance(UserRepositoryInstance repository.UserRepository) *UserServiceImpl {

	return &UserServiceImpl{
		UserRepository: UserRepositoryInstance,
	}
}
func (userService *UserServiceImpl) RegisterUser(request *types.RegisterUserRequest) *types.APIError {

	newUser := data.CreateNewUser(request)

	err := userService.UserRepository.CreateNewUser(newUser)

	if err != nil {

		return &types.APIError{Error: "error while creating user.", ErrorStatus: 500}
	}

	return nil
}

func (userService *UserServiceImpl) LoginUser(request *types.LoginUserRequest) (string, *types.APIError) {

	// insert credential check logic

	userObjectId, err := userService.UserRepository.GetUserCredentials(request)

	if err != nil {

		return "", &types.APIError{Error: "invalid user credentials, please check password.", ErrorStatus: 401}
	}

	JwtToken, ApiError := util.GenerateJWTToken(userObjectId)

	if ApiError != nil {

		return "", ApiError

	}

	return JwtToken, nil
}

func (userService *UserServiceImpl) SendChatRequest(userObjectId primitive.ObjectID, request *types.ChatRequest) (string, *types.APIError) {

	fmt.Println("entered the service method.")

	_, err := userService.UserRepository.GetUserByUserObjectId(userObjectId)

	if err != nil {

		return "", &types.APIError{Error: err.Error(), ErrorStatus: 404}
	} else {

		fmt.Println("user with object id " + userObjectId.Hex() + " was found.")
	}

	_, err = userService.UserRepository.GetUserByUserObjectId(request.ReceiverId)

	if err != nil {

		return "", &types.APIError{Error: err.Error(), ErrorStatus: 404}
	} else {

		fmt.Println("user with object id " + request.ReceiverId.Hex() + " was found.")
	}

	NewChatRequestObjectId, ApiError := userService.UserRepository.SendChatRequest(request)

	if ApiError != nil {

		return "", ApiError
	}

	return NewChatRequestObjectId, nil
}

func (userService *UserServiceImpl) GetOutgoingChatRequestList(string) []types.ChatRequest {

	return nil
}
