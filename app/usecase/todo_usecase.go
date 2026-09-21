package usecase

import (
	"errors"
	"fmt"

	"app/domain/model"
	"app/domain/repository"
	apperrors "app/errors"
)

// ErrTodoNotFound は Todo が存在しない場合と、他人の Todo だった場合の両方で返す。
// 応答を変えると id を総当たりして他人の Todo の実在を判別できてしまうため。
var ErrTodoNotFound = errors.New("Todoが存在しません")

// インターフェースは頭大文字
type TodoUseCase interface {
	GetTodos(userId uint) (*[]model.Todo, error)
	AddTodo(userId uint, title string) error
	DoneTodo(userId uint, id uint) error
	DeleteTodo(userId uint, id uint) error
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

// ownedTodo は自分の Todo だけを返す。存在しない場合も他人の Todo の場合も
// 同じ ErrTodoNotFound を返し、DB エラーだけはそのまま返す。
func (u *todoUseCase) ownedTodo(userId uint, id uint) (*model.Todo, error) {
	todo, err := u.TodoRepository.GetById(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}
	if todo == nil || todo.UserId != userId {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

func (u *todoUseCase) DoneTodo(userId uint, id uint) error {
	todo, err := u.ownedTodo(userId, id)
	if err != nil {
		return err
	}

	todo.Completed = true
	_, err = u.TodoRepository.Update(todo)
	if err != nil {
		return err
	}
	return nil
}

func (u *todoUseCase) DeleteTodo(userId uint, id uint) error {
	todo, err := u.ownedTodo(userId, id)
	if err != nil {
		return err
	}

	err = u.TodoRepository.Delete(todo)
	if err != nil {
		return err
	}
	return nil
}
