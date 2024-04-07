package main

import (
	"context"
	"log"
	//"github.com/gorilla/websocket"
	// "github.com/joho/godotenv"
)

func main() {

	storage, err := NewMongoDBInstance()

	if err != nil {

		log.Fatal(err.Error())
	}

	defer func() {
		if err := storage.db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	APIServer := NewAPIServer(":3000", storage)

	APIServer.Run()
}
