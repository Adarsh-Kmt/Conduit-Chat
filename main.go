package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Adarsh-Kmt/chatapp/controller"
	"github.com/Adarsh-Kmt/chatapp/database"
	"github.com/Adarsh-Kmt/chatapp/repository"
	"github.com/Adarsh-Kmt/chatapp/service"
	"github.com/gorilla/mux"
)

func main() {

	MongoDB, err := database.NewMongoDBInstance()

	if err != nil {

		log.Fatal(err.Error())
	}

	//RedisDB := database.NewRedisDBInstance()

	//chatRepository := repository.NewChatRepositoryImplInstance(RedisDB)
	userRepository := repository.NewUserRepositoryImplInstance(MongoDB)

	userService := service.NewUserServiceImplInstance(userRepository)
	messageService := service.MakeMessageServiceImplInstance(userRepository)

	//chatService := service.NewChatServiceImplInstance(chatRepository)

	userController := controller.NewUserControllerInstance(userService, messageService)

	MainRouter := mux.NewRouter()
	userController.Run(MainRouter)

	http.ListenAndServe(":3000", MainRouter)

	defer func() {
		if err := MongoDB.MongoDBClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

}
