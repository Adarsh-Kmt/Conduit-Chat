package main

import (
	"context"
	//"crypto/x509"
	"fmt"
	"log"
	"time"

	//"github.com/Adarsh-Kmt/chatapp/util"

	"github.com/Adarsh-Kmt/chatapp/types"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

type store interface {
	RegisterUser(*types.User) (string, *APIError)
	SendChatRequest(*types.ChatRequest) (string, *APIError)
	GetChatRequests(string) ([]types.ChatRequest, error)
}

type MongoDBStore struct {
	db *mongo.Client
}

func NewMongoDBInstance() (*MongoDBStore, error) {

	MongoDBInstance := new(MongoDBStore)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	MongoDBInstance.db = client

	err = MongoDBInstance.init()

	if err != nil {

		return nil, err
	}

	return MongoDBInstance, nil
}

func (ms *MongoDBStore) init() error {

	CollectionNames := []string{"ChatUser", "ChatRequest", "Chat"}

	for i := range CollectionNames {

		err := ms.db.Database("chatapp").CreateCollection(context.TODO(), CollectionNames[i])

		if err != nil {
			return fmt.Errorf("error occured while initializing database.")
		}
	}

	collectionNames, err := ms.db.Database("chatapp").ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	// Printing the list of collection names
	fmt.Println("Collections:")
	for _, name := range collectionNames {
		fmt.Println(name)
	}

	return nil
}

func (ms *MongoDBStore) RegisterUser(NewUser *types.User) (string, *types.APIError) {

	//UserCollection := ms.db.Database("chatapp").Collection("ChatUser")

	// ObjectId := NewUserDocument.InsertedID.(primitive.ObjectID)
	// ObjectIdString := ObjectId.Hex()

	// JWTToken, ApiError := util.GenerateJWTToken(ObjectIdString)

	// if ApiError != nil {

	// 	return "", ApiError
	// }
	// return JWTToken, nil

	return "", nil
}

func (ms *MongoDBStore) SendChatRequest(request *types.ChatRequest) (string, *APIError) {

	ChatUserCollection := ms.db.Database("chatapp").Collection("ChatUser")

	var RecieverUser types.User

	err := ChatUserCollection.FindOne(context.TODO(), bson.M{"_id": request.ReceiverId}).Decode(&RecieverUser)

	if err == mongo.ErrNoDocuments {

		return "", &APIError{Error: "no user with id " + request.ReceiverId.Hex(), ErrorStatus: 404}
	}

	var SenderUser types.User

	err = ChatUserCollection.FindOne(context.TODO(), bson.M{"_id": request.ReceiverId}).Decode(&SenderUser)

	if err == mongo.ErrNoDocuments {

		return "", &APIError{Error: "no user with id " + request.ReceiverId.Hex(), ErrorStatus: 404}
	}

	request.SentOnDate = time.Now()

	ChatRequestCollection := ms.db.Database("chatapp").Collection("ChatRequest")

	NewChatRequestDocument, err := ChatRequestCollection.InsertOne(context.TODO(), bson.D{
		{Key: "SenderId", Value: SenderUser.ID},
		{Key: "ReceiverId", Value: RecieverUser.ID},
		{Key: "SentOnDate", Value: request.SentOnDate},
	})

	update := bson.M{
		"$set": bson.M{
			"IncomingChatRequestList." + request.SenderId.Hex(): NewChatRequestDocument.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	filter := bson.M{
		"_id": request.ReceiverId,
	}

	_, err = ChatUserCollection.UpdateOne(context.TODO(), filter, update)

	update = bson.M{
		"$set": bson.M{
			"OutgoingChatRequestList." + request.ReceiverId.Hex(): NewChatRequestDocument.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	filter = bson.M{
		"_id": request.SenderId,
	}

	_, err = ChatUserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {

		return "", &APIError{Error: "error while registering chat request.", ErrorStatus: 500}
	}

	return NewChatRequestDocument.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (ms *MongoDBStore) GetChatRequests(UserId string) ([]types.ChatRequest, error) {

	return nil, nil
}
