package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todoApp/database"
	"todoApp/handlers"
	"todoApp/middleware"

	"github.com/gorilla/mux"
)

func main(){
	//Initialize database
	database.InitDB()

	//Routes
	r := mux.NewRouter()

	//Public routes
	r.HandleFunc("/register", handlers.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)
	
	//Protected routes
	protectRoutes := r.PathPrefix("/").Subrouter()
	protectRoutes.Use(middleware.AuthMiddleware)

	protectRoutes.HandleFunc("/todos", handlers.CreateTodo).Methods(http.MethodPost)
	protectRoutes.HandleFunc("/todos", handlers.GetAllTodos).Methods(http.MethodGet)
	protectRoutes.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods(http.MethodPut)
	protectRoutes.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods(http.MethodDelete)

	//Server
	//Get the PORT from the environment (Render sets this)
    port := os.Getenv("PORT")
    //Run locally if port is empty
    if port == "" {
        port = "8080"
    }
    //Start the server
    fmt.Printf("Server is running on port %s\n", port)

    err := http.ListenAndServe(":"+port, r)
    if err != nil {
        log.Fatal("Server failed to start: ", err)
    }
}