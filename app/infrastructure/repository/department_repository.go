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

func (r *departmentRepository) GetById(id uint) (*model.Department, error) {
	department := &model.Department{}
	err := r.Conn.First(department, id).Error
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (r *departmentRepository) GetList(companyId uint) ([]model.Department, error) {
	var departments []model.Department
	query := r.Conn.Where("")
	// 構造体でクエリを実行する場合、GORM はゼロ以外のフィールドでのみクエリを実行します。つまり、フィールドの値が0、''、falseまたはその他のゼロ値である場合、そのフィールドはクエリ条件の構築に使用されません
	// query = query.Where(model.Department{CompanyId: companyId}, "CompanyId", "Depth").Preload("Ancestors").Preload("Descendants", "ancestor_id != descendant_id").Preload("Descendants.Department")
	// query = query.Where(model.Department{CompanyId: companyId}).Preload("Ancestors", "ancestor_id != descendant_id").Preload("Descendants", "ancestor_id != descendant_id")
	// err := query.Find(&departments).Error
	query = query.Raw("SELECT id, name, depth FROM departments WHERE company_id = ? ORDER BY order_no", companyId)
	err := query.Scan(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *departmentRepository) Create(department *model.Department, parentId *uint) (*model.Department, error) {
	r.Conn.Transaction(func(tx *gorm.DB) error {
		// 仮のソートNoの設定(createするためにソートNoの最大値を設定する)
		var maxOrderNo int
		if err := tx.Table("departments").
			Where(model.Department{CompanyId: department.CompanyId}).
			Select("COALESCE(MAX(order_no), 0) AS maxOrderNo").
			Row().Scan(&maxOrderNo); err != nil {
			fmt.Println(err)
			return err
		}
		department.OrderNo = maxOrderNo + 1
		if err := tx.Create(department).Error; err != nil {
			return err
		}
		// 親がいるかどうか
		// 親がいない場合はdepartment作成して、自己参照を作成して終わり
		// 親がいる場合は親を子孫にもつデータの集合を取得してその子孫として作成したdepartmentを指定する(自己参照忘れない)
		if parentId == nil {
			fmt.Println("親なし")

			descendants := []*model.DepartmentPath{
				{AncestorId: department.ID, DescendantId: department.ID},
			}
			if err := tx.Create(descendants).Error; err != nil {
				return err
			}
		} else {
			fmt.Println("親あり")

			descendants := []model.DepartmentPath{
				{AncestorId: department.ID, DescendantId: department.ID},
			}
			var departmentPaths []model.DepartmentPath
			result := tx.Where(model.DepartmentPath{DescendantId: *parentId}).Select("AncestorId").Find(&departmentPaths)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return gorm.ErrRecordNotFound
			}
			for _, v := range departmentPaths {
				descendants = append(descendants, model.DepartmentPath{AncestorId: v.AncestorId, DescendantId: department.ID})
			}
			// 親のIDから正しいソートNoを設定
			if err := tx.Table("departments").
				Joins("INNER JOIN department_paths AS dp ON departments.id = dp.ancestor_id").
				Joins("INNER JOIN departments AS descendant ON dp.descendant_id = descendant.id").
				Where(model.Department{CompanyId: department.CompanyId}).
				Where("departments.id = ?", parentId).
				Where("departments.depth IN (?)", []int{len(departmentPaths), len(departmentPaths) - 1}). //自分と親の階層を指定
				Select("COALESCE(MAX(descendant.order_no), 0) AS maxOrderNo").
				Row().Scan(&maxOrderNo); err != nil {
				fmt.Println(err)
				return err
			}
			// 取得したOrderNo以降のデータのOrderNoを+1更新する
			if err := tx.Table("departments").
				Where(model.Department{CompanyId: department.CompanyId}).
				Where("departments.order_no > ?", maxOrderNo).
				Where("departments.id != ?", department.ID).
				Update("order_no", gorm.Expr("order_no + 1")).Error; err != nil {
				fmt.Println(err)
				return err
			}
			department.OrderNo = maxOrderNo + 1
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

func (r *departmentRepository) Update(department *model.Department) (*model.Department, error) {
	if err := r.Conn.Save(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (r *departmentRepository) Delete(department *model.Department) error {
	if err := r.Conn.Delete(department).Error; err != nil {
		return err
	}
	return nil
}
