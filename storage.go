package main

import (
	"context"
	//"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

type store interface {
	RegisterUser(*User) (string, *APIError)
	SendChatRequest(*ChatRequest) (string, *APIError)
	GetChatRequests(string) ([]ChatRequest, error)
}

type MongoDBStore struct {
	db *mongo.Client
}

func NewMongoDBInstance() (*MongoDBStore, error) {

	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }

	//uri := os.Getenv("MONGODB_URI")
	// if uri == "" {
	// 	log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	// }
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

func (ms *MongoDBStore) RegisterUser(NewUser *User) (string, *APIError) {

	UserCollection := ms.db.Database("chatapp").Collection("ChatUser")

	NewUserDocument, err := UserCollection.InsertOne(context.TODO(), bson.D{
		{Key: "Name", Value: NewUser.Name},
		{Key: "ChatIdList", Value: NewUser.ChatIdList},
		{Key: "UserConnList", Value: NewUser.UserConnList},
		{Key: "PendingChatRequestList", Value: NewUser.PendingChatRequestList},
	})

	if err != nil {

		return "", &APIError{Error: "error occured while creating user document in the database.", ErrorStatus: 500}
	}

	ObjectId := NewUserDocument.InsertedID.(primitive.ObjectID)
	ObjectIdString := ObjectId.Hex()

	JWTToken, ApiError := GenerateJWTToken(ObjectIdString)

	if ApiError != nil {

		return "", ApiError
	}
	return JWTToken, nil
}

func (ms *MongoDBStore) SendChatRequest(request *ChatRequest) (string, *APIError) {

	// check if the 2 users exist.
	// create a chat request object if the 2 people exist.
	// add chatrequestid and username key:value pair to list for receiver.

	ChatUserCollection := ms.db.Database("chatapp").Collection("ChatUser")

	var recieverUser User

	err := ChatUserCollection.FindOne(context.TODO(), bson.M{"_id": request.ReceiverId}).Decode(&recieverUser)

	if err == mongo.ErrNoDocuments {

		return "", &APIError{Error: "no user with id " + request.ReceiverId.Hex(), ErrorStatus: 404}
	}

	var SenderUser User

	err = ChatUserCollection.FindOne(context.TODO(), bson.M{"_id": request.ReceiverId}).Decode(&SenderUser)

	if err == mongo.ErrNoDocuments {

		return "", &APIError{Error: "no user with id " + request.ReceiverId.Hex(), ErrorStatus: 404}
	}

	request.SentOnDate = time.Now()

	ChatRequestCollection := ms.db.Database("chatapp").Collection("ChatRequest")

	NewChatRequestDocument, err := ChatRequestCollection.InsertOne(context.TODO(), bson.D{
		{Key: "SenderId", Value: SenderUser.ID},
		{Key: "ReceiverId", Value: recieverUser.ID},
		{Key: "SentOnDate", Value: request.SentOnDate},
	})

	update := bson.M{
		"$set": bson.M{
			"PendingChatRequestList." + request.SenderId.Hex(): NewChatRequestDocument.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	filter := bson.M{
		"_id": request.ReceiverId,
	}

	_, err = ChatUserCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {

		return "", &APIError{Error: "error while registering chat request.", ErrorStatus: 500}
	}

	return NewChatRequestDocument.InsertedID.(primitive.ObjectID).Hex(), nil

}

func GenerateJWTToken(UserObjectId string) (string, *APIError) {

	// Parse the RSA private key
	//pkcs1 format
	// pemData, err := os.ReadFile("private_key.pem")
	// if err != nil {
	// 	fmt.Println("Error reading PEM file:", err)
	// 	os.Exit(1)
	// }
	// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)

	secretKey := os.Getenv("JWT_PRIVATE_KEY")
	// if err != nil {
	// 	log.Printf(err.Error())
	// 	return "", &APIError{Error: "error while parsing private key.", ErrorStatus: 500}
	// }

	claims := &jwt.RegisteredClaims{
		Subject:  UserObjectId,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	log.Println(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Println(err.Error())
		return "", &APIError{Error: "error while generating jwt token.", ErrorStatus: 500}
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (*jwt.Token, error) {

	// publicKeyBytes, err := os.ReadFile("public_key.pem")
	// if err != nil {
	// 	return nil, fmt.Errorf("error while reading public key: %v", err)
	// }

	// PublicKeyNew, err := x509.Pa(publicKeyBytes)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	return nil, err
	// }
	// // Parse the PEM-encoded public key
	// block, _ := pem.Decode(publicKeyBytes)
	// if block == nil {
	// 	return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	// }
	// if block.Type != "PUBLIC KEY" {
	// 	return nil, fmt.Errorf("unexpected PEM type: %s", block.Type)
	// }

	// // Parse the public key
	// publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing public key: %v", err)
	// }

	// if _, ok := publicKey.(*rsa.PublicKey); !ok {

	// 	return nil, fmt.Errorf("wrong public key.")
	// }
	log.Println("the token is")
	log.Println(tokenString)

	secretKey := os.Getenv("JWT_PRIVATE_KEY")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("wrong signing key.")
		}

		return []byte(secretKey), nil
	})

}

func getUserObjectIdFromJWT(JWTToken string) (string, error) {

	secretKey := os.Getenv("JWT_PRIVATE_KEY")

	parsedToken, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("wrong signing key.")
		}

		return []byte(secretKey), nil
	})

	if err != nil {

		return "", fmt.Errorf("error in parsing JWT")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("error in getting claims from JWT.")
	}

	userObjectID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("error in extracting subject claim from JWT.")

	}

	return userObjectID, nil

}
func (ms *MongoDBStore) GetChatRequests(UserId string) ([]ChatRequest, error) {

	return nil, nil
}
