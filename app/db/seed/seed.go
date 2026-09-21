package seed

import (
	"fmt"
	"time"

	"app/domain/model"

	"golang.org/x/crypto/bcrypt"
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

// companies は作成順に会社を返す。
// テナント分離が正しく効いているかを確認できるよう、seed は2社分のデータを作る。
func companies(db *gorm.DB) ([]model.Company, error) {
	var list []model.Company
	if err := db.Order("id").Find(&list).Error; err != nil {
		return nil, err
	}
	if len(list) < 2 {
		return nil, fmt.Errorf("会社が2社以上必要です (got %d)", len(list))
	}
	return list, nil
}

// departmentsByName は指定した会社の部署を名前で引けるようにして返す。
func departmentsByName(db *gorm.DB, companyId uint) (map[string]model.Department, error) {
	var list []model.Department
	if err := db.Where("company_id = ?", companyId).Find(&list).Error; err != nil {
		return nil, err
	}
	byName := make(map[string]model.Department, len(list))
	for _, d := range list {
		byName[d.Name] = d
	}
	return byName, nil
}

// selfAndChildPaths は「自分への経路」と「親→子の経路」をまとめて作る。
func selfAndChildPaths(depts map[string]model.Department, parentToChildren map[string][]string) []model.DepartmentPath {
	var paths []model.DepartmentPath
	for _, d := range depts {
		paths = append(paths, model.DepartmentPath{AncestorId: d.ID, DescendantId: d.ID})
	}
	for parent, children := range parentToChildren {
		for _, child := range children {
			paths = append(paths, model.DepartmentPath{AncestorId: depts[parent].ID, DescendantId: depts[child].ID})
		}
	}
	return paths
}

func seedDepartments(db *gorm.DB) error {
	cs, err := companies(db)
	if err != nil {
		return err
	}
	primary, secondary := cs[0], cs[1]

	departments := []model.Department{
		{CompanyId: primary.ID, Name: "経営企画部", Depth: 0, OrderNo: 1},
		{CompanyId: primary.ID, Name: "営業部", Depth: 0, OrderNo: 2},
		{CompanyId: primary.ID, Name: "開発部", Depth: 0, OrderNo: 3},
		{CompanyId: primary.ID, Name: "人事部", Depth: 0, OrderNo: 4},
		{CompanyId: primary.ID, Name: "第一営業課", Depth: 1, OrderNo: 1},
		{CompanyId: primary.ID, Name: "第二営業課", Depth: 1, OrderNo: 2},
		// 2社目。テナント分離の確認用
		{CompanyId: secondary.ID, Name: "総務部", Depth: 0, OrderNo: 1},
		{CompanyId: secondary.ID, Name: "製造部", Depth: 0, OrderNo: 2},
		{CompanyId: secondary.ID, Name: "第一製造課", Depth: 1, OrderNo: 1},
	}
	return db.CreateInBatches(departments, len(departments)).Error
}

func seedDepartmentPaths(db *gorm.DB) error {
	cs, err := companies(db)
	if err != nil {
		return err
	}

	// 会社ごとに「親 -> 配下の課」を定義する。
	// 経路は会社をまたがないので、必ず同じ会社の部署だけで組み立てる。
	hierarchies := []struct {
		companyId uint
		children  map[string][]string
	}{
		{cs[0].ID, map[string][]string{"営業部": {"第一営業課", "第二営業課"}}},
		{cs[1].ID, map[string][]string{"製造部": {"第一製造課"}}},
	}

	var paths []model.DepartmentPath
	for _, h := range hierarchies {
		depts, err := departmentsByName(db, h.companyId)
		if err != nil {
			return err
		}
		if len(depts) == 0 {
			continue
		}
		paths = append(paths, selfAndChildPaths(depts, h.children)...)
	}
	if len(paths) == 0 {
		return nil
	}
	return db.CreateInBatches(paths, len(paths)).Error
}

func seedUsers(db *gorm.DB) error {
	cs, err := companies(db)
	if err != nil {
		return err
	}

	type userData struct {
		companyId uint
		email     string
		password  string
	}
	// company_id は必ず設定する。0 のままだと会社に所属しないユーザーになってしまう。
	entries := []userData{
		{cs[0].ID, "admin@example.com", "password1"},
		{cs[0].ID, "user1@example.com", "password2"},
		{cs[0].ID, "user2@example.com", "password3"},
		// 2社目のユーザー。このアカウントからは1社目のデータが見えないことを確認できる
		{cs[1].ID, "other@example.com", "password4"},
	}

	for _, e := range entries {
		hashed, err := bcrypt.GenerateFromPassword([]byte(e.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := db.Create(&model.User{
			CompanyId: e.companyId,
			Email:     e.email,
			Password:  string(hashed),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedUserInfos(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	type infoData struct {
		lastName  string
		firstName string
		gender    int
		birthDay  time.Time
		working   bool
	}
	data := []infoData{
		{"管理", "太郎", 1, time.Date(1985, 4, 1, 0, 0, 0, 0, time.Local), true},
		{"山田", "花子", 2, time.Date(1992, 8, 15, 0, 0, 0, 0, time.Local), true},
		{"田村", "次郎", 1, time.Date(1990, 12, 3, 0, 0, 0, 0, time.Local), false},
		{"大野", "三郎", 1, time.Date(1988, 6, 20, 0, 0, 0, 0, time.Local), true},
	}

	userInfos := make([]model.UserInfo, len(users))
	for i, u := range users {
		d := data[i%len(data)]
		birthDay := d.birthDay
		userInfos[i] = model.UserInfo{
			UserId:    u.ID,
			LastName:  d.lastName,
			FirstName: d.firstName,
			Gender:    d.gender,
			BirthDay:  &birthDay,
			Working:   d.working,
		}
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

	userAddresses := make([]model.UserAddress, len(users))
	for i, u := range users {
		userAddresses[i] = model.UserAddress{UserId: u.ID}
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

	titles := [][]string{
		{"月次レポート作成", "採用面談の準備"},
		{"営業資料の更新", "顧客へのフォローアップ"},
		{"コードレビュー"},
		{"棚卸しの確認"},
	}

	var todos []model.Todo
	for i, u := range users {
		for j, title := range titles[i%len(titles)] {
			todos = append(todos, model.Todo{UserId: u.ID, Title: title, Completed: j%2 == 0})
		}
	}
	return db.CreateInBatches(todos, len(todos)).Error
}

func seedEmployees(db *gorm.DB) error {
	cs, err := companies(db)
	if err != nil {
		return err
	}
	primary, secondary := cs[0], cs[1]

	employees := []model.Employee{
		{CompanyId: primary.ID, LastName: "田中", FirstName: "太郎", LastNameKana: "タナカ", FirstNameKana: "タロウ", Email: "taro.tanaka@example.com", StaffCode: "EMP001"},
		{CompanyId: primary.ID, LastName: "佐藤", FirstName: "花子", LastNameKana: "サトウ", FirstNameKana: "ハナコ", Email: "hanako.sato@example.com", StaffCode: "EMP002"},
		{CompanyId: primary.ID, LastName: "鈴木", FirstName: "一郎", LastNameKana: "スズキ", FirstNameKana: "イチロウ", Email: "ichiro.suzuki@example.com", StaffCode: "EMP003"},
		{CompanyId: primary.ID, LastName: "高橋", FirstName: "美咲", LastNameKana: "タカハシ", FirstNameKana: "ミサキ", Email: "misaki.takahashi@example.com", StaffCode: "EMP004"},
		{CompanyId: primary.ID, LastName: "伊藤", FirstName: "健太", LastNameKana: "イトウ", FirstNameKana: "ケンタ", Email: "kenta.ito@example.com", StaffCode: "EMP005"},
		// 2社目。1社目のアカウントからは見えないことを確認できる
		{CompanyId: secondary.ID, LastName: "山本", FirstName: "健", LastNameKana: "ヤマモト", FirstNameKana: "ケン", Email: "ken.yamamoto@example.com", StaffCode: "EMP101"},
		{CompanyId: secondary.ID, LastName: "中村", FirstName: "由美", LastNameKana: "ナカムラ", FirstNameKana: "ユミ", Email: "yumi.nakamura@example.com", StaffCode: "EMP102"},
	}
	return db.CreateInBatches(employees, len(employees)).Error
}

func seedEmployeeDepartments(db *gorm.DB) error {
	var employees []model.Employee
	if err := db.Order("id").Find(&employees).Error; err != nil {
		return err
	}
	if len(employees) == 0 {
		return nil
	}

	// 部署は会社ごとに引く。会社をまたいだ紐付けを作らないため。
	deptsByCompany := make(map[uint][]model.Department)

	var records []model.EmployeeDepartments
	for i, emp := range employees {
		depts, ok := deptsByCompany[emp.CompanyId]
		if !ok {
			if err := db.Where("company_id = ?", emp.CompanyId).Order("order_no").Find(&depts).Error; err != nil {
				return err
			}
			deptsByCompany[emp.CompanyId] = depts
		}
		if len(depts) == 0 {
			continue
		}
		records = append(records, model.EmployeeDepartments{
			CompanyId:    emp.CompanyId,
			EmployeeId:   emp.ID,
			DepartmentId: depts[i%len(depts)].ID,
		})
	}
	if len(records) == 0 {
		return nil
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
