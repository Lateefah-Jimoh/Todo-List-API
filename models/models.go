package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct{
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Todos []Todo `json:"todos"` //foreign key
}

type Todo struct{
	gorm.Model
	Title string `json:"title"`
	Description string  `json:"description"`
	Completed bool `json:"completed"`
	CreatedAt time.Time	`json:"created_at"`
	UserID uint `json:"user_id"` //foreign key
}
