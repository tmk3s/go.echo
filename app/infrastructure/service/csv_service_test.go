package service_test

import (
	"bytes"
	"encoding/csv"
	"mime/multipart"
	"testing"
	"time"

	"app/domain/model"
	svc "app/infrastructure/service"
)

// nopFile wraps bytes.Reader to satisfy multipart.File
type nopFile struct{ *bytes.Reader }

func (nopFile) Close() error { return nil }

func newFile(content string) multipart.File {
	return nopFile{bytes.NewReader([]byte(content))}
}

func newFileWithBOM(content string) multipart.File {
	bom := []byte{0xEF, 0xBB, 0xBF}
	return nopFile{bytes.NewReader(append(bom, []byte(content)...))}
}

// --- ParseEmployeeRows ---

func TestParseEmployeeRows_BasicFields(t *testing.T) {
	content := "スタッフコード,姓,名,姓（カナ）,名（カナ）,メールアドレス,郵便番号,都道府県,市区町村,住所1,住所2,電話番号,入社日,退職日,退職区分,ステータス\n" +
		"EMP001,田中,太郎,タナカ,タロウ,taro@example.com,100-0001,東京都,千代田区,千代田1-1,,03-1234-5678,,,,"

	rows, err := svc.NewCsvService().ParseEmployeeRows(newFile(content))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	r := rows[0]
	checks := map[string][2]string{
		"StaffCode":      {r.StaffCode, "EMP001"},
		"LastName":       {r.LastName, "田中"},
		"FirstName":      {r.FirstName, "太郎"},
		"LastNameKana":   {r.LastNameKana, "タナカ"},
		"FirstNameKana":  {r.FirstNameKana, "タロウ"},
		"Email":          {r.Email, "taro@example.com"},
		"PostCode":       {r.PostCode, "100-0001"},
		"PrefectureName": {r.PrefectureName, "東京都"},
		"City":           {r.City, "千代田区"},
		"Tel":            {r.Tel, "03-1234-5678"},
	}
	for field, pair := range checks {
		if pair[0] != pair[1] {
			t.Errorf("%s: want %q, got %q", field, pair[1], pair[0])
		}
	}
}

func TestParseEmployeeRows_BOMStripped(t *testing.T) {
	content := "スタッフコード,メールアドレス\nEMP001,emp@example.com"

	rows, err := svc.NewCsvService().ParseEmployeeRows(newFileWithBOM(content))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %d", len(rows))
	}
	if rows[0].StaffCode != "EMP001" {
		t.Errorf("BOM not stripped: want EMP001, got %q", rows[0].StaffCode)
	}
}

func TestParseEmployeeRows_DepartmentColumns(t *testing.T) {
	content := "スタッフコード,メールアドレス,部署1,部署2\n" +
		"EMP001,emp@example.com,営業部,第一営業課"

	rows, err := svc.NewCsvService().ParseEmployeeRows(newFile(content))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows[0].DepartmentNames) != 2 {
		t.Fatalf("want 2 departments, got %d", len(rows[0].DepartmentNames))
	}
	if rows[0].DepartmentNames[0] != "営業部" {
		t.Errorf("dept[0]: want 営業部, got %s", rows[0].DepartmentNames[0])
	}
	if rows[0].DepartmentNames[1] != "第一営業課" {
		t.Errorf("dept[1]: want 第一営業課, got %s", rows[0].DepartmentNames[1])
	}
}

func TestParseEmployeeRows_Tenure(t *testing.T) {
	content := "スタッフコード,メールアドレス,入社日,退職日,退職区分,ステータス\n" +
		"EMP001,emp@example.com,2020-04-01,2023-03-31,自己都合,退職済"

	rows, err := svc.NewCsvService().ParseEmployeeRows(newFile(content))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows[0].Tenures) != 1 {
		t.Fatalf("want 1 tenure, got %d", len(rows[0].Tenures))
	}
	ten := rows[0].Tenures[0]
	if ten.JoinedOn != "2020-04-01" {
		t.Errorf("JoinedOn: want 2020-04-01, got %s", ten.JoinedOn)
	}
	if ten.ResignationOn != "2023-03-31" {
		t.Errorf("ResignationOn: want 2023-03-31, got %s", ten.ResignationOn)
	}
	if ten.ResignationType != "自己都合" {
		t.Errorf("ResignationType: want 自己都合, got %s", ten.ResignationType)
	}
	if ten.Status != "退職済" {
		t.Errorf("Status: want 退職済, got %s", ten.Status)
	}
}

func TestParseEmployeeRows_NoTenureWhenJoinedOnEmpty(t *testing.T) {
	content := "スタッフコード,メールアドレス,入社日\nEMP001,emp@example.com,"

	rows, err := svc.NewCsvService().ParseEmployeeRows(newFile(content))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows[0].Tenures) != 0 {
		t.Errorf("want 0 tenures when 入社日 is empty, got %d", len(rows[0].Tenures))
	}
}

func TestParseEmployeeRows_SkipsEmptyRows(t *testing.T) {
	// スタッフコードもメールも空の行はスキップ
	content := "スタッフコード,メールアドレス\n" +
		"EMP001,emp1@example.com\n" +
		",\n" +
		"EMP002,emp2@example.com\n"

	rows, err := svc.NewCsvService().ParseEmployeeRows(newFile(content))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Errorf("want 2 rows (empty row skipped), got %d", len(rows))
	}
}

// --- GenerateEmployeeCSV ---

func TestGenerateEmployeeCSV_BOM(t *testing.T) {
	data, err := svc.NewCsvService().GenerateEmployeeCSV([]model.Employee{})
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 3 || data[0] != 0xEF || data[1] != 0xBB || data[2] != 0xBF {
		t.Error("output does not start with UTF-8 BOM")
	}
}

func TestGenerateEmployeeCSV_BaseHeaders(t *testing.T) {
	data, err := svc.NewCsvService().GenerateEmployeeCSV([]model.Employee{})
	if err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(data[3:])) // skip BOM
	headers, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"スタッフコード", "姓", "名", "姓（カナ）", "名（カナ）", "メールアドレス",
		"郵便番号", "都道府県", "市区町村", "住所1", "住所2", "電話番号",
		"入社日", "退職日", "退職区分", "ステータス",
	}
	if len(headers) != len(want) {
		t.Fatalf("want %d headers, got %d: %v", len(want), len(headers), headers)
	}
	for i, h := range want {
		if headers[i] != h {
			t.Errorf("header[%d]: want %q, got %q", i, h, headers[i])
		}
	}
}

func TestGenerateEmployeeCSV_DynamicDepartmentHeaders(t *testing.T) {
	employees := []model.Employee{
		{
			StaffCode:   "EMP001",
			Departments: []model.Department{{Name: "営業部"}, {Name: "第一営業課"}},
		},
	}

	data, err := svc.NewCsvService().GenerateEmployeeCSV(employees)
	if err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(data[3:]))
	headers, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}

	if headers[len(headers)-2] != "部署1" {
		t.Errorf("want 部署1, got %q", headers[len(headers)-2])
	}
	if headers[len(headers)-1] != "部署2" {
		t.Errorf("want 部署2, got %q", headers[len(headers)-1])
	}
}

func TestGenerateEmployeeCSV_EmployeeData(t *testing.T) {
	joinedOn := time.Date(2020, 4, 1, 0, 0, 0, 0, time.Local)
	resignedOn := time.Date(2023, 3, 31, 0, 0, 0, 0, time.Local)
	resignationType := "自己都合"
	status := "退職済"
	postCode := "100-0001"
	prefName := "東京都"
	city := "千代田区"
	addr1 := "千代田1-1"

	employees := []model.Employee{
		{
			StaffCode:     "EMP001",
			LastName:      "田中",
			FirstName:     "太郎",
			LastNameKana:  "タナカ",
			FirstNameKana: "タロウ",
			Email:         "taro@example.com",
			Address: &model.EmployeeAddress{
				PostCode:       &postCode,
				PrefectureName: &prefName,
				City:           &city,
				AddressLine1:   &addr1,
			},
			Tenures: []model.EmployeeTenures{
				{
					JoinedOn:        joinedOn,
					ResignationOn:   &resignedOn,
					ResignationType: &resignationType,
					Status:          &status,
				},
			},
		},
	}

	data, err := svc.NewCsvService().GenerateEmployeeCSV(employees)
	if err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(data[3:]))
	if _, err := r.Read(); err != nil { // skip header
		t.Fatal(err)
	}
	record, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}

	// col index: 0=code,1=last,2=first,3=lastKana,4=firstKana,5=email,6=post,7=pref,8=city,9=addr1,10=addr2,11=tel,12=joined,13=resigned,14=type,15=status
	wantByIdx := map[int]string{
		0:  "EMP001",
		1:  "田中",
		2:  "太郎",
		3:  "タナカ",
		4:  "タロウ",
		5:  "taro@example.com",
		6:  "100-0001",
		7:  "東京都",
		8:  "千代田区",
		9:  "千代田1-1",
		12: "2020-04-01",
		13: "2023-03-31",
		14: "自己都合",
		15: "退職済",
	}
	for idx, want := range wantByIdx {
		if record[idx] != want {
			t.Errorf("col[%d]: want %q, got %q", idx, want, record[idx])
		}
	}
}

func TestGenerateEmployeeCSV_EmptyAddressAndTenure(t *testing.T) {
	employees := []model.Employee{
		{StaffCode: "EMP001", Email: "emp@example.com"},
	}

	data, err := svc.NewCsvService().GenerateEmployeeCSV(employees)
	if err != nil {
		t.Fatal(err)
	}

	r := csv.NewReader(bytes.NewReader(data[3:]))
	r.Read() // skip header
	record, _ := r.Read()

	// address and tenure fields should be empty strings
	for _, idx := range []int{6, 7, 8, 9, 10, 11, 12, 13, 14, 15} {
		if record[idx] != "" {
			t.Errorf("col[%d] should be empty, got %q", idx, record[idx])
		}
	}
}
