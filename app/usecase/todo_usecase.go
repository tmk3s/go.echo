package usecase

import (
	"fmt"

	"app/domain/model"
	"app/domain/repository"
)

// インターフェースは頭大文字
type TodoUseCase interface {
	GetTodos(userId uint) (*[]model.Todo, error)
	AddTodo(userId uint, title string) error
	DoneTodo(id uint) error
	DeleteTodo(id uint) error
}

type todoUseCase struct {
	repository.TodoRepository
	repository.UserRepository
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

func (u *todoUseCase) AddTodo(userId uint, title string) error {
	user, err := u.UserRepository.GetById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("ユーザーが存在しません")
	}

	todo := model.NewTodo(userId, title)
	_, err = u.TodoRepository.Add(todo)
	if err != nil {
		return err
	}
	return nil
}

func (u *todoUseCase) DoneTodo(id uint) error {
	todo, err := u.TodoRepository.GetById(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("Todoが存在しません")
	}

	todo.Completed = true
	_, err = u.TodoRepository.Update(todo)
	if err != nil {
		return err
	}
	return nil
}

func (u *todoUseCase) DeleteTodo(id uint) error {
	todo, err := u.TodoRepository.GetById(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("Todoが存在しません")
	}

	err = u.TodoRepository.Delete(todo)
	if err != nil {
		return err
	}
	return nil
}
