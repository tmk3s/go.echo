package repository

import (
	"gorm.io/gorm"
	"app/domain/model"
	"app/domain/repository"
)

type todoRepository struct {
	Conn *gorm.DB
}

func NewTodoRepository(Conn *gorm.DB) repository.TodoRepository {
	return &todoRepository{Conn} // ポインタを返す
}

func (r *todoRepository) GetList(userId uint) ([]model.Todo, error) {
	var todos []model.Todo
	query := r.Conn.Where("")
	query = query.Where(model.Todo{UserId: userId})
	err := query.Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, err
}

func (r *todoRepository) Add(u *model.Todo) (*model.Todo, error) {
	var todo model.Todo
	return &todo, nil
}

func (r *todoRepository) Done(u *model.Todo) (*model.Todo, error) {
	var todo model.Todo
	return &todo, nil
}

func (r *todoRepository) Delete(id uint) (error) {
	return nil
}