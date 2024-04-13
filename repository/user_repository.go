package repository

import (
	"fmt"

	"github.com/Adarsh-Kmt/chatapp/data"
	"github.com/Adarsh-Kmt/chatapp/database"
	"github.com/Adarsh-Kmt/chatapp/types"

	"context"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	CreateNewUser(*types.User) *types.APIError
	SendChatRequest(*types.ChatRequest) (string, *types.APIError)
	GetIncomingChatRequestList(primitive.ObjectID) ([]types.ChatRequest, *types.APIError)
	GetOutgoingChatRequestList(primitive.ObjectID) ([]types.ChatRequest, *types.APIError)
	GetUserCredentials(*types.LoginUserRequest) (string, error)
	GetUserByUserObjectId(primitive.ObjectID) (*types.User, error)
}

type UserRepositoryImpl struct {
	Mongo *database.MongoDBStore
}

func NewUserRepositoryImplInstance(MongoDBInstance *database.MongoDBStore) *UserRepositoryImpl {

	return &UserRepositoryImpl{
		Mongo: MongoDBInstance,
	}

}
func (userRepository *UserRepositoryImpl) CreateNewUser(NewUser *types.User) *types.APIError {

	UserCollection := userRepository.Mongo.DB.Database("chatapp").Collection("ChatUser")

	_, err := UserCollection.InsertOne(context.TODO(), data.CreateUserBSON(NewUser))

	if err != nil {

		return &types.APIError{Error: "error occured while creating user document in the database.", ErrorStatus: 500}
	}

	return nil

}

func (UserRepository *UserRepositoryImpl) GetUserCredentials(loginRequest *types.LoginUserRequest) (string, error) {

	UserCollection := UserRepository.Mongo.DB.Database("chatapp").Collection("ChatUser")

	var registeredUser types.User
	err := UserCollection.FindOne(context.TODO(), bson.M{"Name": loginRequest.Name, "Password": loginRequest.Password}).Decode(&registeredUser)

	if err == mongo.ErrNoDocuments {

		return "", fmt.Errorf("unauthorized")
	}

	return registeredUser.ID.Hex(), nil

}

func (userRepository *UserRepositoryImpl) GetOutgoingChatRequestList(userObjectId primitive.ObjectID) ([]types.ChatRequest, *types.APIError) {

	UserCollection := userRepository.Mongo.DB.Database("chatapp").Collection("ChatUser")
	projection := bson.M{"OutgoingChatRequestList": 1}

	filter := bson.M{"_id": userObjectId}

	var outgoingChatRequestListResponse types.OutgoingChatRequestListResponse

	err := UserCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&outgoingChatRequestListResponse)

	if err == mongo.ErrNoDocuments {

		return nil, &types.APIError{Error: "No outgoing chat requests.", ErrorStatus: 404}
	}

	return outgoingChatRequestListResponse.OutgoingChatRequestList, nil
}
func (userRepository *UserRepositoryImpl) GetIncomingChatRequestList(userObjectId primitive.ObjectID) ([]types.ChatRequest, *types.APIError) {

	UserCollection := userRepository.Mongo.DB.Database("chatapp").Collection("ChatUser")

	filter := bson.M{"_id": userObjectId}

	projection := bson.M{"IncomingChatRequestList": 1}

	var incomingChatRequestListResponse types.IncomingChatRequestListResponse

	err := UserCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&incomingChatRequestListResponse)

	if err == mongo.ErrNoDocuments {

		return nil, &types.APIError{Error: "no incoming chat requests.", ErrorStatus: 404}
	}

	return incomingChatRequestListResponse.IncomingChatRequestList, nil
}

func (userRepository *UserRepositoryImpl) SendChatRequest(chatRequest *types.ChatRequest) (string, *types.APIError) {

	fmt.Println("entered repository.")
	UserCollection := userRepository.Mongo.DB.Database("chatapp").Collection("ChatUser")

	SenderChatRequestObjectIdString, SenderChatRequestBSON := data.CreateChatRequestBSON(chatRequest)
	_, ReceiverChatRequestBSON := data.CreateChatRequestBSON(chatRequest)

	SenderUpdateResult, err := UserCollection.UpdateOne(context.TODO(), bson.M{"_id": chatRequest.SenderId}, bson.M{"$push": bson.M{"OutgoingChatRequestList": SenderChatRequestBSON}})

	if SenderUpdateResult.MatchedCount == 0 {

		return "", &types.APIError{Error: "user with object id " + chatRequest.SenderId.Hex() + " does not exist.", ErrorStatus: 404}
	}
	if err != nil {

		return "", &types.APIError{Error: "error while registering chat request.", ErrorStatus: 500}
	} else {

		fmt.Println("updated sender user.")
	}
	ReceiverUpdateResult, err := UserCollection.UpdateOne(context.TODO(), bson.M{"_id": chatRequest.ReceiverId}, bson.M{"$push": bson.M{"IncomingChatRequestList": ReceiverChatRequestBSON}})

	if err != nil {

		return "", &types.APIError{Error: "error while registering chat request.", ErrorStatus: 500}
	} else {

		fmt.Println("updated receiver user.")
	}

	if ReceiverUpdateResult.MatchedCount == 0 {

		return "", &types.APIError{Error: "user with object id " + chatRequest.ReceiverId.Hex() + " does not exist.", ErrorStatus: 404}
	}

	return SenderChatRequestObjectIdString, nil
}

func (userRepository *UserRepositoryImpl) GetUserByUserObjectId(userObjectId primitive.ObjectID) (*types.User, error) {

	UserCollection := userRepository.Mongo.DB.Database("chatapp").Collection("ChatUser")

	var User types.User

	err := UserCollection.FindOne(context.TODO(), bson.M{"_id": userObjectId}).Decode(&User)

	if err == mongo.ErrNoDocuments {

		return nil, fmt.Errorf("user with object id " + userObjectId.Hex() + " not found.")
	}

	return &User, nil
}
