package repository

import {
	"fmt"

	"gor.io/gorm"
	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
}

type todoRepository struct {
	Conn *gorm.DB
}

func NewTodoRepository(conn *gorm.DB) repository.TodoRepository {
	return &todoRepository(conn) // ポインタを返す
}

func (r *todoRepository) GetList(userId uint) ([]model.Todo, error) {
	var todos []model.Todo
	query = r.Conn.Where("")
	query = query.Where(model.Todo{UserId: userId})
	err != query.Find(&todos).Error
	if err != nil {
		return nil, echo.ErrNotFound
	}
	return todos, err
}

func (r *todoRepository) Add(u *model.Todo) (*model.Todo, error) {

}

func (r *todoRepository) Done(u *model.Todo) (*model.Todo, error) {

}

func (r *todoRepository) Delete(id uint) (error) {

}