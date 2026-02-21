package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todoApp/database"
	"todoApp/middleware"
	"todoApp/models"
	"todoApp/utils"

	"github.com/gorilla/mux"
)

//Register
func Register(w http.ResponseWriter, r *http.Request) {
	//decode request body
	var reg *models.User
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//check if the user already exists
	var user *models.User
	err = database.Db.Where("email= ?", reg.Email).First(&user).Error
	if err == nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	//Hash the password
	HashPassword, err := utils.HashPassword(reg.Password)
	if err != nil {
		http.Error(w, "unable to hash password", http.StatusBadRequest)
		return
	}

	reg.Password = HashPassword

	// add the user to the database
	err = database.Db.Create(&reg).Error
	if err != nil {
		http.Error(w, "unable to create user", http.StatusBadRequest)
		return
	}

	// send a response
	fmt.Println("User created successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully",})
}

//Login 
func Login(w http.ResponseWriter, r *http.Request) {
	// decode the request the request body
	var login models.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if the user exists
	var user models.User
	err = database.Db.Where("email = ?", login.Email).First(&user).Error
	if err != nil {
		http.Error(w, "this user does not exist", http.StatusBadRequest)
		return
	}

	// check if password matches what we have in our database
	err = utils.ComparePassword(login.Password, user.Password)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	// uint ---> int ---> string.

	idStr := strconv.Itoa(int(user.ID))

	// generating a token
	token, err := middleware.GenerateJWT(idStr)
	if err != nil {
		http.Error(w, "unable to generate token", http.StatusInternalServerError)
		return
	}

	// send a response
	fmt.Println("Login successfully")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

//Create Todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	//check if the user is logged in
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//decode the request body
	var task models.Todo
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Assign the request to the user
	task.UserID = userID

	//Create the task in database
	err = database.Db.Create(&task).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//send request response
	fmt.Println("Todo created successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo created successfully",})

}

//Get All Todos
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	//Check if the user is logged in
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Create slice to hold the response
	var todos []models.Todo
	//To get all tasks owned by the user using the userID
	err = database.Db.Where("user_id = ?", userID).Find(&todos).Error
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
	//Send response
	fmt.Println("Todo generated successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(todos)
}

//Update Todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	//Check if the user is logged in
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//get todoID using path parameter
	vars := mux.Vars(r)
	todoID := vars["id"]
	if todoID == "" {
		http.Error(w, "Todo ID is required", http.StatusBadRequest)
		return
	}
	//Check if the todo exist in the database
	var todo models.Todo
	err = database.Db.First(&todo, todoID).Error
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}
	//Check if the user is the owner of the todo in the database
	if todo.UserID != userID {
		http.Error(w, "You are not authorized to edit this todo", http.StatusForbidden)
		return
	}
	//decode the request
	var updateData struct {
		Title string `json:"title"`
		Description string  `json:"description"`
		Completed bool `json:"completed"` 
	}
	err = json.NewDecoder(r.Body).Decode(&updateData) 
	if err != nil {
		http.Error(w, "Invalid input JSON", http.StatusBadRequest)
		return
	}
	//update and save
	todo.Title = updateData.Title
	todo.Description = updateData.Description
	todo.Completed = updateData.Completed

	err = database.Db.Save(&todo).Error
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	fmt.Println("Todo updated successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully",})
	
}

//Delete Todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	//Check if the user is logged in
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//get todoID using path parameter
	vars := mux.Vars(r)
	todoID := vars["id"]
	if todoID == "" {
		http.Error(w, "Todo ID is required", http.StatusBadRequest)
		return
	}
	//Check if the todo exist in the database
	var todo models.Todo
	err = database.Db.First(&todo, todoID).Error
	if err != nil {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}
	//Check if the user is the owner of the todo in the database
	if todo.UserID != userID {
		http.Error(w, "You are not authorized to delete this todo", http.StatusForbidden)
		return
	}
	//delete todo from the database
	err = database.Db.Delete(&todo).Error
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	fmt.Println("Todo deleted successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully",})

}
