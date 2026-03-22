package service

import (
	"encoding/csv"
	"io"
	"mime/multipart"

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
