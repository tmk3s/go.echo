package repository

import (
	"errors"
	"fmt"

	"app/domain/model"
	"app/domain/repository"
	"gorm.io/gorm"
)

type departmentRepository struct {
	Conn *gorm.DB
}

func NewDepartmentRepository(Conn *gorm.DB) repository.DepartmentRepository {
	return &departmentRepository{Conn}
}

func (r *departmentRepository) GetList(companyId uint) ([]model.Department, error) {
	var departments []model.Department
	query := r.Conn.Where("")
	query = query.Where(model.Department{CompanyId: companyId})
	err := query.Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, err
}

func (r *departmentRepository) Create(department *model.Department, parentId *uint) (*model.Department, error) {
	if err := r.Conn.Create(department).Error; err != nil {
		return nil, err
	}

	// 親がいるかどうか
	// 親がいない場合はdepartment作成して、自己参照を作成して終わり
	// 親がいる場合は親を子孫にもつデータの集合を取得してその子孫として作成したdepartmentを指定する(自己参照忘れない)
	if parentId == nil {
		descendants := []*model.DepartmentPath{
			{Ancestor: department.ID, Descendant: department.ID},
		}
		if err := r.Conn.Create(descendants).Error; err != nil {
			return nil, err
		}
	} else {
		fmt.Println("OK!2")
		descendants := []model.DepartmentPath{
			{Ancestor: department.ID, Descendant: department.ID},
		}
		var departmentPaths []model.DepartmentPath
		result := r.Conn.Where(model.DepartmentPath{Descendant: *parentId}).Select("Ancestor").Find(&departmentPaths)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return department, nil
		}
		for _, v := range departmentPaths {
			descendants = append(descendants, model.DepartmentPath{Ancestor: v.Ancestor, Descendant: department.ID})
		}
		if err := r.Conn.Create(&descendants).Error; err != nil {
			return nil, err
		}
	}

	return department, nil
}

// func (r *departmentRepository) Update(department *model.Department) (*model.Department, error) {
// 	if err := r.Conn.Save(department).Error; err != nil {
// 		return nil, err
// 	}
// 	return department, nil
// }

// func (r *departmentRepository) Delete(department *model.Department) error {
// 	if err := r.Conn.Delete(department).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
