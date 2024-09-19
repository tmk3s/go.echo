package model

// import "fmt"
import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	UserId    uint   `json:"user_id" gorm:"praimaly_key"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func NewTodo(userId uint, title string) *Todo {
	return &Todo{
		UserId:    userId,
		Title:     title,
		Completed: false,
	}
}
