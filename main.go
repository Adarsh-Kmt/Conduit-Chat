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

	storage, err := database.NewMongoDBInstance()

	if err != nil {

		log.Fatal(err.Error())
	}

	userRepository := repository.NewUserRepositoryImplInstance(storage)

	userService := service.NewUserServiceImplInstance(userRepository)

	userController := controller.NewUserControllerInstance(userService)
	MainRouter := mux.NewRouter()
	userController.Run(MainRouter)

	http.ListenAndServe(":3000", MainRouter)

	defer func() {
		if err := storage.DB.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

}
