package repository

import "app/domain/model"

// 実装はapp/infrastructure/repository/todo_repository.goにて行う
type TodoRepository interface {
	GetById(id uint) (*model.Todo, error)
	GetList(id uint) ([]model.Todo, error)
	Add(todo *model.Todo) (*model.Todo, error)
	Update(todo *model.Todo) (*model.Todo, error)
	Delete(todo *model.Todo) error
}
