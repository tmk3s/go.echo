package repository

import (
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

func (r *departmentRepository) GetById(companyId uint, id uint) (*model.Department, error) {
	department := &model.Department{}
	// company_id を条件に含めないと他社の部署を取得できてしまう
	if err := r.Conn.Where("company_id = ?", companyId).First(department, id).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (r *departmentRepository) GetList(companyId uint) ([]model.Department, error) {
	var departments []model.Department
	// 構造体でクエリを実行すると GORM はゼロ値のフィールドを条件に使わないため、
	// companyId が 0 のときに条件ごと消えて全社のデータが返ってしまう。明示的に条件を書く。
	err := r.Conn.Where("company_id = ?", companyId).
		Select("id, name, depth").
		Order("order_no").
		Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *departmentRepository) Create(department *model.Department, parentId *uint) (*model.Department, error) {
	err := r.Conn.Transaction(func(tx *gorm.DB) error {
		// 仮のソートNoの設定(createするためにソートNoの最大値を設定する)
		var maxOrderNo int
		if err := tx.Table("departments").
			Where("company_id = ?", department.CompanyId).
			Select("COALESCE(MAX(order_no), 0) AS maxOrderNo").
			Row().Scan(&maxOrderNo); err != nil {
			return err
		}
		department.OrderNo = maxOrderNo + 1
		if err := tx.Create(department).Error; err != nil {
			return err
		}

		// 親がいない場合はdepartment作成して、自己参照を作成して終わり
		if parentId == nil {
			descendants := []*model.DepartmentPath{
				{AncestorId: department.ID, DescendantId: department.ID},
			}
			return tx.Create(descendants).Error
		}

		// 親が自社の部署かを検証する(他社の部署を親に指定させない)
		var parent model.Department
		if err := tx.Where("company_id = ?", department.CompanyId).First(&parent, *parentId).Error; err != nil {
			return err
		}

		// 親がいる場合は親を子孫にもつデータの集合を取得して
		// その子孫として作成したdepartmentを指定する(自己参照忘れない)
		descendants := []model.DepartmentPath{
			{AncestorId: department.ID, DescendantId: department.ID},
		}
		var departmentPaths []model.DepartmentPath
		// Find はレコードが0件でも ErrRecordNotFound を返さないため、件数で判定する
		if err := tx.Where("descendant_id = ?", *parentId).
			Select("ancestor_id").
			Find(&departmentPaths).Error; err != nil {
			return err
		}
		if len(departmentPaths) == 0 {
			return gorm.ErrRecordNotFound
		}
		for _, v := range departmentPaths {
			descendants = append(descendants, model.DepartmentPath{AncestorId: v.AncestorId, DescendantId: department.ID})
		}
		// 親のIDから正しいソートNoを設定
		if err := tx.Table("departments").
			Joins("INNER JOIN department_paths AS dp ON departments.id = dp.ancestor_id").
			Joins("INNER JOIN departments AS descendant ON dp.descendant_id = descendant.id").
			Where("departments.company_id = ?", department.CompanyId).
			Where("departments.id = ?", parentId).
			Where("departments.depth IN (?)", []int{len(departmentPaths), len(departmentPaths) - 1}). //自分と親の階層を指定
			Select("COALESCE(MAX(descendant.order_no), 0) AS maxOrderNo").
			Row().Scan(&maxOrderNo); err != nil {
			return err
		}
		// 取得したOrderNo以降のデータのOrderNoを+1更新する(自社のデータのみ)
		if err := tx.Table("departments").
			Where("departments.company_id = ?", department.CompanyId).
			Where("departments.order_no > ?", maxOrderNo).
			Where("departments.id != ?", department.ID).
			Update("order_no", gorm.Expr("order_no + 1")).Error; err != nil {
			return err
		}
		department.OrderNo = maxOrderNo + 1
		department.Depth = len(departmentPaths)
		if err := tx.Save(department).Error; err != nil {
			return err
		}
		return tx.Create(&descendants).Error
	})
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (r *departmentRepository) Update(department *model.Department) (*model.Department, error) {
	if err := r.Conn.Save(department).Error; err != nil {
		return nil, err
	}
	return department, nil
}

func (r *departmentRepository) Delete(department *model.Department) error {
	// 子孫を取得
	descendantPaths := []model.DepartmentPath{}
	if err := r.Conn.Where("ancestor_id = ?", department.ID).Find(&descendantPaths).Error; err != nil {
		return err
	}
	// 先祖を取得
	ancestorPaths := []model.DepartmentPath{}
	if err := r.Conn.Where("descendant_id = ?", department.ID).Find(&ancestorPaths).Error; err != nil {
		return err
	}
	// 子孫のマスタを取得(自社のデータのみ)
	descendants := []model.Department{}
	var descendantIds []uint
	for _, descendant := range descendantPaths {
		descendantIds = append(descendantIds, descendant.DescendantId)
	}
	if len(descendantIds) > 0 {
		if err := r.Conn.Where("company_id = ?", department.CompanyId).
			Where("id IN (?)", descendantIds).
			Find(&descendants).Error; err != nil {
			return err
		}
	}

	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if len(ancestorPaths) > 0 {
			if err := tx.Delete(&ancestorPaths).Error; err != nil {
				return err
			}
		}
		if len(descendantPaths) > 0 {
			if err := tx.Delete(&descendantPaths).Error; err != nil {
				return err
			}
		}
		if len(descendants) > 0 {
			if err := tx.Delete(&descendants).Error; err != nil {
				return err
			}
		}
		return tx.Delete(department).Error
	})
}
