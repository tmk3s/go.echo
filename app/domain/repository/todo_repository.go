package repository

import "app/domain/model"

// 実装はapp/infrastructure/repository/todo_repository.goにて行う
type TodoRepository interface {
	GetList(id uint) ([]model.Todo, error)
	Add(u *model.Todo) (*model.Todo, error)
	Done(u *model.Todo) (*model.Todo, error)
	Delete(id uint) error
}