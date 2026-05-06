package service

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	domainservice "app/domain/service"

	"github.com/jszwec/csvutil"
)

var reDept = regexp.MustCompile(`^部署(\d+)$`)

var baseJPHeaders = []string{
	"スタッフコード", "姓", "名", "姓（カナ）", "名（カナ）", "メールアドレス",
	"郵便番号", "都道府県", "市区町村", "住所1", "住所2", "電話番号",
	"入社日", "退職日", "退職区分", "ステータス",
}

type csvService struct{}

type departmentRow struct {
	Name string `csv:"name"`
}

func NewCsvService() domainservice.CsvService {
	return &csvService{}
}

func (s *csvService) ParseDepartmentNames(file multipart.File) ([]string, error) {
	reader := csv.NewReader(file)
	dec, err := csvutil.NewDecoder(reader)
	if err != nil {
		return nil, err
	}

	var names []string
	for {
		var row departmentRow
		if err := dec.Decode(&row); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if row.Name != "" {
			names = append(names, row.Name)
		}
	}
	return names, nil
}

func (s *csvService) GenerateEmployeeCSV(data []domainservice.EmployeeExportData) ([]byte, error) {
	maxDepts := 0
	for _, d := range data {
		if len(d.Departments) > maxDepts {
			maxDepts = len(d.Departments)
		}
	}

	headers := append([]string{}, baseJPHeaders...)
	for i := 1; i <= maxDepts; i++ {
		headers = append(headers, fmt.Sprintf("部署%d", i))
	}

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF}) // BOM for Excel
	w := csv.NewWriter(&buf)
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, d := range data {
		row := []string{
			d.Employee.StaffCode,
			d.Employee.LastName,
			d.Employee.FirstName,
			d.Employee.LastNameKana,
			d.Employee.FirstNameKana,
			d.Employee.Email,
		}

		var postCode, prefName, city, addr1, addr2, tel string
		if d.Address != nil {
			postCode = derefStr(d.Address.PostCode)
			prefName = derefStr(d.Address.PrefectureName)
			city = derefStr(d.Address.City)
			addr1 = derefStr(d.Address.AddressLine1)
			addr2 = derefStr(d.Address.AddressLine2)
			tel = derefStr(d.Address.Tel)
		}
		row = append(row, postCode, prefName, city, addr1, addr2, tel)

		// 在籍情報は先頭の1件のみ出力
		var joinedOn, resignationOn, resignationType, status string
		if len(d.Tenures) > 0 {
			t := d.Tenures[0]
			joinedOn = t.JoinedOn.Format("2006-01-02")
			resignationOn = formatTimePtr(t.ResignationOn)
			resignationType = derefStr(t.ResignationType)
			status = derefStr(t.Status)
		}
		row = append(row, joinedOn, resignationOn, resignationType, status)

		for i := 0; i < maxDepts; i++ {
			if i < len(d.Departments) {
				row = append(row, d.Departments[i].Name)
			} else {
				row = append(row, "")
			}
		}

		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), nil
}

func (s *csvService) ParseEmployeeRows(file multipart.File) ([]domainservice.EmployeeCSVRow, error) {
	r := csv.NewReader(skipBOM(file))

	rawHeaders, err := r.Read()
	if err != nil {
		return nil, err
	}

	baseIdx := make(map[string]int)
	deptIdx := make(map[int]int)
	maxDeptN := 0

	for i, h := range rawHeaders {
		h = strings.TrimSpace(h)
		if m := reDept.FindStringSubmatch(h); m != nil {
			n, _ := strconv.Atoi(m[1])
			deptIdx[n] = i
			if n > maxDeptN {
				maxDeptN = n
			}
		} else {
			baseIdx[h] = i
		}
	}

	get := func(record []string, key string) string {
		if i, ok := baseIdx[key]; ok && i < len(record) {
			return strings.TrimSpace(record[i])
		}
		return ""
	}

	var rows []domainservice.EmployeeCSVRow
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		staffCode := get(record, "スタッフコード")
		email := get(record, "メールアドレス")
		if staffCode == "" && email == "" {
			continue
		}

		row := domainservice.EmployeeCSVRow{
			StaffCode:      staffCode,
			LastName:       get(record, "姓"),
			FirstName:      get(record, "名"),
			LastNameKana:   get(record, "姓（カナ）"),
			FirstNameKana:  get(record, "名（カナ）"),
			Email:          email,
			PostCode:       get(record, "郵便番号"),
			PrefectureName: get(record, "都道府県"),
			City:           get(record, "市区町村"),
			AddressLine1:   get(record, "住所1"),
			AddressLine2:   get(record, "住所2"),
			Tel:            get(record, "電話番号"),
		}

		// 在籍情報は1件のみ
		if joinedOn := get(record, "入社日"); joinedOn != "" {
			row.Tenures = append(row.Tenures, domainservice.TenureCSVRow{
				JoinedOn:        joinedOn,
				ResignationOn:   get(record, "退職日"),
				ResignationType: get(record, "退職区分"),
				Status:          get(record, "ステータス"),
			})
		}

		for n := 1; n <= maxDeptN; n++ {
			if i, ok := deptIdx[n]; ok && i < len(record) {
				if name := strings.TrimSpace(record[i]); name != "" {
					row.DepartmentNames = append(row.DepartmentNames, name)
				}
			}
		}

		rows = append(rows, row)
	}
	return rows, nil
}

func skipBOM(file multipart.File) io.Reader {
	buf := bufio.NewReader(file)
	bom, err := buf.Peek(3)
	if err == nil && bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		buf.Discard(3) //nolint:errcheck
	}
	return buf
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
