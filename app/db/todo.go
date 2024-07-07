package db

import "fmt"

import (
	"gorm.io/gorm"
)

type Todo struct {
  gorm.Model
	UserId     uint   `json:"user_id" gorm:"praimaly_key"`
  Title      string `json:"title"`
  Completed  bool   `json:"completed"`
}

type Todos []Todo

func CreateTodo(todo *Todo) {
  db.Create(todo)
}

func FindTodo(id string) Todo {
	var todo Todo
	db.Where("id = ?", id).First(&todo)
	return todo
}

func FindTodos(t *Todo) Todos {
	var todos Todos
	db.Where(t).Find(&todos)
	fmt.Print(todos)
	return todos
}

func DeleteTodo(t *Todo) error {
	if rows := db.Delete(&t).RowsAffected; rows == 0 {
		return fmt.Errorf("Could not find Todo (%v) to delete!", t)
	}
	return nil
}

func UpdateTodo(t *Todo) error {
	t.Completed = true
	db.Save(&t)
	return nil
}