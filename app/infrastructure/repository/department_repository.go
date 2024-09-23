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
	// 構造体でクエリを実行する場合、GORM はゼロ以外のフィールドでのみクエリを実行します。つまり、フィールドの値が0、''、falseまたはその他のゼロ値である場合、そのフィールドはクエリ条件の構築に使用されません
	// query = query.Where(model.Department{CompanyId: companyId}, "CompanyId", "Depth").Preload("Ancestors").Preload("Descendants", "ancestor != descendant").Preload("Descendants.Department")
	// query = query.Where(model.Department{CompanyId: companyId}).Preload("Ancestors", "ancestor != descendant").Preload("Descendants", "ancestor != descendant")
	// err := query.Find(&departments).Error
	query = query.Raw("SELECT d_child.id, d_child.name, d_child.depth FROM departments AS d INNER JOIN department_paths AS dp ON d.id = dp.ancestor INNER JOIN departments AS d_child ON dp.descendant = d_child.id WHERE d.company_id = ? AND d.depth = 0 ORDER BY dp.ancestor, d_child.depth, d_child.id", companyId)
	err := query.Scan(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *departmentRepository) Create(department *model.Department, parentId *uint) (*model.Department, error) {
	r.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(department).Error; err != nil {
			return err
		}
		// 親がいるかどうか
		// 親がいない場合はdepartment作成して、自己参照を作成して終わり
		// 親がいる場合は親を子孫にもつデータの集合を取得してその子孫として作成したdepartmentを指定する(自己参照忘れない)
		if parentId == nil {
			fmt.Println("親なし")
	
			descendants := []*model.DepartmentPath{
				{Ancestor: department.ID, Descendant: department.ID},
			}
			if err := tx.Create(descendants).Error; err != nil {
				return err
			}
		} else {
			fmt.Println("親あり")
	
			descendants := []model.DepartmentPath{
				{Ancestor: department.ID, Descendant: department.ID},
			}
			var departmentPaths []model.DepartmentPath
			result := tx.Where(model.DepartmentPath{Descendant: *parentId}).Select("Ancestor").Find(&departmentPaths)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			for _, v := range departmentPaths {
				descendants = append(descendants, model.DepartmentPath{Ancestor: v.Ancestor, Descendant: department.ID})
			}
			department.Depth = len(departmentPaths)
			if err := tx.Save(department).Error; err != nil {
				return err
			}
	
			if err := tx.Create(&descendants).Error; err != nil {
				return err
			}
		}
		return nil
	})
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
