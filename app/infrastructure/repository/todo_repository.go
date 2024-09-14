package repository

import (
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
	err := r.Conn.First(todo, id).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
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
