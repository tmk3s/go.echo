package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"

	domainservice "app/domain/service"

	"github.com/jszwec/csvutil"
)

type csvService struct {
	rows [][]string
}

func NewCsvService() domainservice.CsvService {
	return &csvService{}
}

func (s *csvService) ConvertCsv(file multipart.File, fileHeader *multipart.FileHeader) error {
	reader := csv.NewReader(file)
	dec, err := csvutil.NewDecoder(reader)
	if err != nil {
		return err
	}

	for {
		var row []string
		if err := dec.Decode(&row); err == io.EOF {
			break
		} else if err != nil {
			if e, ok := err.(*csv.ParseError); ok {
				return fmt.Errorf("データ形式に誤りがあります: %v (StartLine: %d, Line: %d, Column: %d)", e.Err, e.StartLine, e.Line, e.Column)
			}
			return err
		}
		s.rows = append(s.rows, row)
	}
	return nil
}
