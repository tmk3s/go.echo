package usecase

import (
	"app/domain/model"
	"app/domain/repository"
	// "app/domain/service"
)

// インターフェースは頭大文字
type TodoUseCase interface {
	GetTodos(userId uint) (*[]model.Todo, error)
	AddTodo(params TodoCreateParams) error
	DoneTodo(id uint) error
	DeleteTodo(id uint) error
}

type todoUseCase struct {
	repository.TodoRepository
	repository.UserRepository
}

type TodoCreateParams struct {
	UserId uint
	Title string
}

func NewTodoUseCase(r repository.TodoRepository, ur repository.UserRepository) TodoUseCase {
	return &todoUseCase{r, ur}
}

func (u *todoUseCase) GetTodos(userId uint) (*[]model.Todo, error) {
	todos, err := u.TodoRepository.GetList(userId)
	if err != nil {
		return nil, err
	}
	return &todos, nil
}

func (u *todoUseCase) AddTodo(params TodoCreateParams) error {
	return nil
}

func (u *todoUseCase) DoneTodo(id uint) error {
	return nil
}

func (u *todoUseCase) DeleteTodo(id uint) error {
	return nil
}

