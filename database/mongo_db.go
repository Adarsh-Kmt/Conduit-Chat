package database

import (
	"context"
	"fmt"
	"log"

	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStore struct {
	DB *mongo.Client
}

func NewMongoDBInstance() (*MongoDBStore, error) {

	MongoDBInstance := new(MongoDBStore)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	MongoDBInstance.DB = client

	err = MongoDBInstance.init()

	if err != nil {

		return nil, err
	}

	return MongoDBInstance, nil
}

func (ms *MongoDBStore) init() error {

	CollectionNames := []string{"ChatUser", "ChatRequest", "Chat"}

	for i := range CollectionNames {

		err := ms.DB.Database("chatapp").CreateCollection(context.TODO(), CollectionNames[i])

		if err != nil {
			return fmt.Errorf("error occured while initializing database.")
		}
	}

	collectionNames, err := ms.DB.Database("chatapp").ListCollectionNames(context.TODO(), bson.D{})
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
