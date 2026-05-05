package seed

import (
	"app/domain/model"
	"time"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	steps := []func(*gorm.DB) error{
		seedTags,
		seedPostCodes,
		seedPrefectures,
		seedCompanies,
		seedDepartments,
		seedDepartmentPaths,
		seedUsers,
		seedUserInfos,
		seedUserAddresses,
		seedTodos,
		seedEmployees,
		seedEmployeeDepartments,
		seedEmployeeAddresses,
		seedEmployeeTenures,
	}
	for _, step := range steps {
		if err := step(db); err != nil {
			return err
		}
	}
	return nil
}

func seedTags(db *gorm.DB) error {
	tags := []model.Tag{
		{Name: "重要"},
		{Name: "緊急"},
		{Name: "確認待ち"},
		{Name: "対応中"},
		{Name: "完了"},
	}
	return db.CreateInBatches(tags, len(tags)).Error
}

func seedPostCodes(db *gorm.DB) error {
	postCodes := []model.PostCode{
		{Code: "100-0001", PrefectureName: "東京都", CityName: "千代田区", TownAreaName: "千代田"},
		{Code: "530-0001", PrefectureName: "大阪府", CityName: "大阪市北区", TownAreaName: "梅田"},
		{Code: "460-0001", PrefectureName: "愛知県", CityName: "名古屋市中区", TownAreaName: "栄"},
		{Code: "220-0001", PrefectureName: "神奈川県", CityName: "横浜市西区", TownAreaName: "高島"},
		{Code: "060-0001", PrefectureName: "北海道", CityName: "札幌市中央区", TownAreaName: "北一条西"},
	}
	return db.CreateInBatches(postCodes, len(postCodes)).Error
}

func seedPrefectures(db *gorm.DB) error {
	prefectures := []model.Prefecture{
		{Name: "北海道"},
		{Name: "青森県"},
		{Name: "岩手県"},
		{Name: "宮城県"},
		{Name: "秋田県"},
		{Name: "山形県"},
		{Name: "福島県"},
		{Name: "茨城県"},
		{Name: "栃木県"},
		{Name: "群馬県"},
		{Name: "埼玉県"},
		{Name: "千葉県"},
		{Name: "東京都"},
		{Name: "神奈川県"},
		{Name: "新潟県"},
		{Name: "富山県"},
		{Name: "石川県"},
		{Name: "福井県"},
		{Name: "山梨県"},
		{Name: "長野県"},
		{Name: "岐阜県"},
		{Name: "静岡県"},
		{Name: "愛知県"},
		{Name: "三重県"},
		{Name: "滋賀県"},
		{Name: "京都府"},
		{Name: "大阪府"},
		{Name: "兵庫県"},
		{Name: "奈良県"},
		{Name: "和歌山県"},
		{Name: "鳥取県"},
		{Name: "島根県"},
		{Name: "岡山県"},
		{Name: "広島県"},
		{Name: "山口県"},
		{Name: "徳島県"},
		{Name: "香川県"},
		{Name: "愛媛県"},
		{Name: "高知県"},
		{Name: "福岡県"},
		{Name: "佐賀県"},
		{Name: "長崎県"},
		{Name: "熊本県"},
		{Name: "大分県"},
		{Name: "宮崎県"},
		{Name: "鹿児島県"},
		{Name: "沖縄県"},
	}
	return db.CreateInBatches(prefectures, len(prefectures)).Error
}

func seedCompanies(db *gorm.DB) error {
	companies := []model.Company{
		{Name: "株式会社サンプル"},
		{Name: "テスト工業株式会社"},
	}
	return db.CreateInBatches(companies, len(companies)).Error
}

func seedDepartments(db *gorm.DB) error {
	var company model.Company
	if err := db.First(&company).Error; err != nil {
		return err
	}

	departments := []model.Department{
		{CompanyId: company.ID, Name: "経営企画部", Depth: 0, OrderNo: 1},
		{CompanyId: company.ID, Name: "営業部", Depth: 0, OrderNo: 2},
		{CompanyId: company.ID, Name: "開発部", Depth: 0, OrderNo: 3},
		{CompanyId: company.ID, Name: "人事部", Depth: 0, OrderNo: 4},
		{CompanyId: company.ID, Name: "第一営業課", Depth: 1, OrderNo: 1},
		{CompanyId: company.ID, Name: "第二営業課", Depth: 1, OrderNo: 2},
	}
	return db.CreateInBatches(departments, len(departments)).Error
}

func seedDepartmentPaths(db *gorm.DB) error {
	var departments []model.Department
	if err := db.Find(&departments).Error; err != nil {
		return err
	}
	if len(departments) < 6 {
		return nil
	}

	// 営業部(index:1) -> 第一営業課(index:4), 第二営業課(index:5)
	salesDept := departments[1]
	subDept1 := departments[4]
	subDept2 := departments[5]

	paths := []model.DepartmentPath{
		// 各部署は自身へのパスを持つ
		{AncestorId: departments[0].ID, DescendantId: departments[0].ID},
		{AncestorId: departments[1].ID, DescendantId: departments[1].ID},
		{AncestorId: departments[2].ID, DescendantId: departments[2].ID},
		{AncestorId: departments[3].ID, DescendantId: departments[3].ID},
		{AncestorId: subDept1.ID, DescendantId: subDept1.ID},
		{AncestorId: subDept2.ID, DescendantId: subDept2.ID},
		// 営業部 -> 配下の課
		{AncestorId: salesDept.ID, DescendantId: subDept1.ID},
		{AncestorId: salesDept.ID, DescendantId: subDept2.ID},
	}
	return db.CreateInBatches(paths, len(paths)).Error
}

func seedUsers(db *gorm.DB) error {
	var company model.Company
	if err := db.First(&company).Error; err != nil {
		return err
	}

	users := []model.User{
		{CompanyId: company.ID, Email: "admin@example.com", Password: "hashed_password_1"},
		{CompanyId: company.ID, Email: "user1@example.com", Password: "hashed_password_2"},
		{CompanyId: company.ID, Email: "user2@example.com", Password: "hashed_password_3"},
	}
	return db.CreateInBatches(users, len(users)).Error
}

func seedUserInfos(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	userInfos := []model.UserInfo{
		{UserId: users[0].ID, LastName: "管理", FirstName: "太郎", Gender: 1, BirthDay: time.Date(1985, 4, 1, 0, 0, 0, 0, time.Local), Working: true},
		{UserId: users[1].ID, LastName: "山田", FirstName: "花子", Gender: 2, BirthDay: time.Date(1992, 8, 15, 0, 0, 0, 0, time.Local), Working: true},
		{UserId: users[2].ID, LastName: "田村", FirstName: "次郎", Gender: 1, BirthDay: time.Date(1990, 12, 3, 0, 0, 0, 0, time.Local), Working: false},
	}
	return db.CreateInBatches(userInfos, len(userInfos)).Error
}

func seedUserAddresses(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	userAddresses := []model.UserAddress{
		{UserId: users[0].ID},
		{UserId: users[1].ID},
		{UserId: users[2].ID},
	}
	return db.CreateInBatches(userAddresses, len(userAddresses)).Error
}

func seedTodos(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	todos := []model.Todo{
		{UserId: users[0].ID, Title: "月次レポート作成", Completed: true},
		{UserId: users[0].ID, Title: "採用面談の準備", Completed: false},
		{UserId: users[1].ID, Title: "営業資料の更新", Completed: false},
		{UserId: users[1].ID, Title: "顧客へのフォローアップ", Completed: true},
		{UserId: users[2].ID, Title: "コードレビュー", Completed: false},
	}
	return db.CreateInBatches(todos, len(todos)).Error
}

func seedEmployees(db *gorm.DB) error {
	var company model.Company
	if err := db.First(&company).Error; err != nil {
		return err
	}

	employees := []model.Employee{
		{CompanyId: company.ID, LastName: "田中", FirstName: "太郎", LastNameKana: "タナカ", FirstNameKana: "タロウ", Email: "taro.tanaka@example.com", StaffCode: "EMP001"},
		{CompanyId: company.ID, LastName: "佐藤", FirstName: "花子", LastNameKana: "サトウ", FirstNameKana: "ハナコ", Email: "hanako.sato@example.com", StaffCode: "EMP002"},
		{CompanyId: company.ID, LastName: "鈴木", FirstName: "一郎", LastNameKana: "スズキ", FirstNameKana: "イチロウ", Email: "ichiro.suzuki@example.com", StaffCode: "EMP003"},
		{CompanyId: company.ID, LastName: "高橋", FirstName: "美咲", LastNameKana: "タカハシ", FirstNameKana: "ミサキ", Email: "misaki.takahashi@example.com", StaffCode: "EMP004"},
		{CompanyId: company.ID, LastName: "伊藤", FirstName: "健太", LastNameKana: "イトウ", FirstNameKana: "ケンタ", Email: "kenta.ito@example.com", StaffCode: "EMP005"},
	}
	return db.CreateInBatches(employees, len(employees)).Error
}

func seedEmployeeDepartments(db *gorm.DB) error {
	var employees []model.Employee
	if err := db.Find(&employees).Error; err != nil {
		return err
	}

	var departments []model.Department
	if err := db.Find(&departments).Error; err != nil {
		return err
	}

	if len(employees) == 0 || len(departments) == 0 {
		return nil
	}

	idx := func(i int) int {
		if i < len(departments) {
			return i
		}
		return len(departments) - 1
	}

	records := []model.EmployeeDepartments{
		{CompanyId: employees[0].CompanyId, EmployeeId: employees[0].ID, DepartmentId: departments[idx(0)].ID},
		{CompanyId: employees[1].CompanyId, EmployeeId: employees[1].ID, DepartmentId: departments[idx(1)].ID},
		{CompanyId: employees[2].CompanyId, EmployeeId: employees[2].ID, DepartmentId: departments[idx(2)].ID},
		{CompanyId: employees[3].CompanyId, EmployeeId: employees[3].ID, DepartmentId: departments[idx(1)].ID},
		{CompanyId: employees[4].CompanyId, EmployeeId: employees[4].ID, DepartmentId: departments[idx(3)].ID},
	}
	return db.CreateInBatches(records, len(records)).Error
}

func seedEmployeeAddresses(db *gorm.DB) error {
	var employees []model.Employee
	if err := db.Find(&employees).Error; err != nil {
		return err
	}
	if len(employees) == 0 {
		return nil
	}

	type addrData struct {
		postCode  string
		prefId    uint
		city      string
		addrLine1 string
		tel       string
	}

	data := []addrData{
		{"100-0001", 13, "千代田区", "千代田1-1-1", "03-1234-5678"},
		{"530-0001", 27, "北区", "梅田1-2-3", "06-1234-5678"},
		{"460-0001", 23, "中区", "栄1-1-1", "052-123-4567"},
		{"220-0001", 14, "西区", "高島1-1-1", "045-123-4567"},
		{"060-0001", 1, "中央区", "北一条西1-1", "011-123-4567"},
	}

	records := make([]model.EmployeeAddress, len(employees))
	for i, emp := range employees {
		d := data[i%len(data)]
		postCode := d.postCode
		prefId := d.prefId
		city := d.city
		addrLine1 := d.addrLine1
		tel := d.tel
		records[i] = model.EmployeeAddress{
			CompanyId:    emp.CompanyId,
			EmployeeId:   emp.ID,
			PostCode:     &postCode,
			PrefectureId: &prefId,
			City:         &city,
			AddressLine1: &addrLine1,
			Tel:          &tel,
		}
	}
	return db.CreateInBatches(records, len(records)).Error
}

func seedEmployeeTenures(db *gorm.DB) error {
	var employees []model.Employee
	if err := db.Find(&employees).Error; err != nil {
		return err
	}
	if len(employees) == 0 {
		return nil
	}

	active := "在職中"
	resigned := "退職済"
	voluntary := "自己都合"
	resignedOn := time.Date(2023, 3, 31, 0, 0, 0, 0, time.Local)

	type tenureData struct {
		joinedOn        time.Time
		resignationOn   *time.Time
		resignationType *string
		status          string
	}

	data := []tenureData{
		{time.Date(2018, 4, 1, 0, 0, 0, 0, time.Local), nil, nil, active},
		{time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local), nil, nil, active},
		{time.Date(2015, 10, 1, 0, 0, 0, 0, time.Local), &resignedOn, &voluntary, resigned},
		{time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local), nil, nil, active},
		{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), nil, nil, active},
	}

	records := make([]model.EmployeeTenures, len(employees))
	for i, emp := range employees {
		d := data[i%len(data)]
		status := d.status
		records[i] = model.EmployeeTenures{
			CompanyId:       emp.CompanyId,
			EmployeeId:      emp.ID,
			JoinedOn:        d.joinedOn,
			ResignationOn:   d.resignationOn,
			ResignationType: d.resignationType,
			Status:          &status,
		}
	}
	return db.CreateInBatches(records, len(records)).Error
}
