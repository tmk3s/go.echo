package repository

import (
	"errors"
	"fmt"

	apperrors "app/errors"

	"app/domain/model"
	"app/domain/repository"
	"gorm.io/gorm"
)

type todoRepository struct {
	Conn *gorm.DB
}

func NewTodoRepository(Conn *gorm.DB) repository.TodoRepository {
	return &todoRepository{Conn} // ポインタを返す
}

func (r *todoRepository) GetById(id uint) (*model.Todo, error) {
	todo := &model.Todo{}
	if err := r.Conn.First(todo, id).Error; err != nil {
		// gorm のエラーをそのまま返すと usecase が gorm に依存するため共通エラーに変換する
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo(id=%d): %w", id, apperrors.ErrNotFound)
		}
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) GetList(userId uint) ([]model.Todo, error) {
	var todos []model.Todo
	// 構造体でのクエリはゼロ値のフィールドが条件から外れるため、明示的に条件を書く
	err := r.Conn.Where("user_id = ?", userId).Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, err
}

func (r *todoRepository) Add(todo *model.Todo) (*model.Todo, error) {
	if err := r.Conn.Create(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	if err := r.Conn.Save(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) Delete(todo *model.Todo) error {
	if err := r.Conn.Delete(todo).Error; err != nil {
		return err
	}
	return nil
}
