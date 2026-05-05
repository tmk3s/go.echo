package service

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	domainservice "app/domain/service"

	"github.com/jszwec/csvutil"
)

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
	maxTenures := 0
	maxDepts := 0
	for _, d := range data {
		if len(d.Tenures) > maxTenures {
			maxTenures = len(d.Tenures)
		}
		if len(d.Departments) > maxDepts {
			maxDepts = len(d.Departments)
		}
	}

	headers := []string{
		"staff_code", "last_name", "first_name", "last_name_kana", "first_name_kana", "email",
		"post_code", "prefecture_name", "city", "address_line1", "address_line2", "tel",
	}
	for i := 1; i <= maxTenures; i++ {
		headers = append(headers,
			fmt.Sprintf("joined_on_%d", i),
			fmt.Sprintf("resignation_on_%d", i),
			fmt.Sprintf("resignation_type_%d", i),
			fmt.Sprintf("status_%d", i),
		)
	}
	for i := 1; i <= maxDepts; i++ {
		headers = append(headers, fmt.Sprintf("department_%d", i))
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

		for i := 0; i < maxTenures; i++ {
			if i < len(d.Tenures) {
				t := d.Tenures[i]
				row = append(row,
					t.JoinedOn.Format("2006-01-02"),
					formatTimePtr(t.ResignationOn),
					derefStr(t.ResignationType),
					derefStr(t.Status),
				)
			} else {
				row = append(row, "", "", "", "")
			}
		}

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

	headers, err := r.Read()
	if err != nil {
		return nil, err
	}

	colIndex := make(map[string]int, len(headers))
	for i, h := range headers {
		colIndex[strings.TrimSpace(h)] = i
	}

	maxTenures := 0
	maxDepts := 0
	for h := range colIndex {
		if strings.HasPrefix(h, "joined_on_") {
			if n, err := strconv.Atoi(strings.TrimPrefix(h, "joined_on_")); err == nil && n > maxTenures {
				maxTenures = n
			}
		}
		if strings.HasPrefix(h, "department_") {
			if n, err := strconv.Atoi(strings.TrimPrefix(h, "department_")); err == nil && n > maxDepts {
				maxDepts = n
			}
		}
	}

	get := func(record []string, key string) string {
		if i, ok := colIndex[key]; ok && i < len(record) {
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

		staffCode := get(record, "staff_code")
		email := get(record, "email")
		if staffCode == "" && email == "" {
			continue
		}

		row := domainservice.EmployeeCSVRow{
			StaffCode:      staffCode,
			LastName:       get(record, "last_name"),
			FirstName:      get(record, "first_name"),
			LastNameKana:   get(record, "last_name_kana"),
			FirstNameKana:  get(record, "first_name_kana"),
			Email:          email,
			PostCode:       get(record, "post_code"),
			PrefectureName: get(record, "prefecture_name"),
			City:           get(record, "city"),
			AddressLine1:   get(record, "address_line1"),
			AddressLine2:   get(record, "address_line2"),
			Tel:            get(record, "tel"),
		}

		for i := 1; i <= maxTenures; i++ {
			joinedOn := get(record, fmt.Sprintf("joined_on_%d", i))
			if joinedOn == "" {
				continue
			}
			row.Tenures = append(row.Tenures, domainservice.TenureCSVRow{
				JoinedOn:        joinedOn,
				ResignationOn:   get(record, fmt.Sprintf("resignation_on_%d", i)),
				ResignationType: get(record, fmt.Sprintf("resignation_type_%d", i)),
				Status:          get(record, fmt.Sprintf("status_%d", i)),
			})
		}

		for i := 1; i <= maxDepts; i++ {
			if name := get(record, fmt.Sprintf("department_%d", i)); name != "" {
				row.DepartmentNames = append(row.DepartmentNames, name)
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
