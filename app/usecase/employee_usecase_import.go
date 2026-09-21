package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
)

// BulkImportEnqueuer is the interface for enqueuing bulk import background tasks.
type BulkImportEnqueuer interface {
	EnqueueBulkCreate(jobId uint) error
	EnqueueBulkUpdate(jobId uint) error
}

func (u *employeeUseCase) EnqueueBulkCreate(companyId uint, file multipart.File) (*model.Job, error) {
	csvBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	job, err := u.jobRepo.Create(&model.Job{
		CompanyId: companyId,
		JobType:   "employee:bulk_create",
		Status:    "pending",
		FileData:  csvBytes,
	})
	if err != nil {
		return nil, err
	}

	if err := u.enqueuer.EnqueueBulkCreate(job.ID); err != nil {
		u.jobRepo.Fail(job.ID, err.Error())
		return nil, err
	}
	return job, nil
}

func (u *employeeUseCase) EnqueueBulkUpdate(companyId uint, file multipart.File) (*model.Job, error) {
	csvBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	job, err := u.jobRepo.Create(&model.Job{
		CompanyId: companyId,
		JobType:   "employee:bulk_update",
		Status:    "pending",
		FileData:  csvBytes,
	})
	if err != nil {
		return nil, err
	}

	if err := u.enqueuer.EnqueueBulkUpdate(job.ID); err != nil {
		u.jobRepo.Fail(job.ID, err.Error())
		return nil, err
	}
	return job, nil
}

// csvReadCloser wraps bytes.Reader to satisfy multipart.File (adds no-op Close).
type csvReadCloser struct{ *bytes.Reader }

func (csvReadCloser) Close() error { return nil }

func (u *employeeUseCase) ProcessBulkCreate(ctx context.Context, jobId uint) error {
	job, err := u.jobRepo.GetById(jobId)
	if err != nil {
		return err
	}

	rows, err := u.csvService.ParseEmployeeRows(csvReadCloser{bytes.NewReader(job.FileData)})
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	lookup, err := u.buildImportLookup(job.CompanyId)
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	// 同じスタッフコードが何行あっても、エラーメッセージには1回だけ載せる
	// (error_message は jobs テーブルに保存されるため、件数分並べると肥大化する)
	counts := make(map[string]int, len(rows))
	codes := make([]string, 0, len(rows))
	for _, row := range rows {
		if counts[row.StaffCode] == 0 {
			codes = append(codes, row.StaffCode)
		}
		counts[row.StaffCode]++
	}

	var duplicates []string
	var inFileDuplicates []string
	for _, code := range codes {
		if _, exists := lookup.empByCode[code]; exists {
			duplicates = append(duplicates, code)
			continue
		}
		// 既存データだけでなく、アップロードされたファイル内の重複も弾く
		if counts[code] > 1 {
			inFileDuplicates = append(inFileDuplicates, code)
		}
	}
	if len(duplicates) > 0 {
		msg := fmt.Sprintf("以下のスタッフコードはすでに登録されています: %s", strings.Join(duplicates, ", "))
		u.jobRepo.Fail(jobId, msg)
		return fmt.Errorf("%s", msg)
	}
	if len(inFileDuplicates) > 0 {
		msg := fmt.Sprintf("アップロードされたファイル内でスタッフコードが重複しています: %s", strings.Join(inFileDuplicates, ", "))
		u.jobRepo.Fail(jobId, msg)
		return fmt.Errorf("%s", msg)
	}

	if err := u.jobRepo.StartProcessing(jobId, len(rows)); err != nil {
		return err
	}

	// 途中で失敗した場合に部分的な行が残ると再アップロードできなくなるため、
	// 取り込み全体を1つのトランザクションにまとめる
	if err := u.repo.Transaction(func(repo repository.EmployeeRepository) error {
		for i, row := range rows {
			select {
			case <-ctx.Done():
				return fmt.Errorf("処理がキャンセルされました: %w", ctx.Err())
			default:
			}

			emp := &model.Employee{
				CompanyId:     job.CompanyId,
				StaffCode:     row.StaffCode,
				LastName:      row.LastName,
				FirstName:     row.FirstName,
				LastNameKana:  row.LastNameKana,
				FirstNameKana: row.FirstNameKana,
				Email:         row.Email,
			}
			created, err := repo.Create(emp)
			if err != nil {
				return fmt.Errorf("行%d: %w", i+1, err)
			}
			if err := u.applyRelated(repo, job.CompanyId, created.ID, row, lookup); err != nil {
				return fmt.Errorf("行%d: %w", i+1, err)
			}
			// 進捗はジョブ用の別コネクションで更新するため、ロールバックの影響を受けない
			u.jobRepo.UpdateProgress(jobId, i+1)
		}
		return nil
	}); err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	return u.jobRepo.Complete(jobId)
}

func (u *employeeUseCase) ProcessBulkUpdate(ctx context.Context, jobId uint) error {
	job, err := u.jobRepo.GetById(jobId)
	if err != nil {
		return err
	}

	rows, err := u.csvService.ParseEmployeeRows(csvReadCloser{bytes.NewReader(job.FileData)})
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	lookup, err := u.buildImportLookup(job.CompanyId)
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	var missing []string
	for _, row := range rows {
		if _, exists := lookup.empByCode[row.StaffCode]; !exists {
			missing = append(missing, row.StaffCode)
		}
	}
	if len(missing) > 0 {
		msg := fmt.Sprintf("以下のスタッフコードに紐づく社員が見つかりません: %s", strings.Join(missing, ", "))
		u.jobRepo.Fail(jobId, msg)
		return fmt.Errorf("%s", msg)
	}

	if err := u.jobRepo.StartProcessing(jobId, len(rows)); err != nil {
		return err
	}

	// 取り込み全体を1つのトランザクションにまとめる(途中失敗時に部分更新を残さない)
	if err := u.repo.Transaction(func(repo repository.EmployeeRepository) error {
		for i, row := range rows {
			select {
			case <-ctx.Done():
				return fmt.Errorf("処理がキャンセルされました: %w", ctx.Err())
			default:
			}

			emp := lookup.empByCode[row.StaffCode]
			emp.LastName = row.LastName
			emp.FirstName = row.FirstName
			emp.LastNameKana = row.LastNameKana
			emp.FirstNameKana = row.FirstNameKana
			emp.Email = row.Email
			if err := repo.UpdateEmployee(emp); err != nil {
				return fmt.Errorf("行%d: %w", i+1, err)
			}
			if err := u.applyRelated(repo, job.CompanyId, emp.ID, row, lookup); err != nil {
				return fmt.Errorf("行%d: %w", i+1, err)
			}
			u.jobRepo.UpdateProgress(jobId, i+1)
		}
		return nil
	}); err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	return u.jobRepo.Complete(jobId)
}

type importLookup struct {
	empByCode  map[string]*model.Employee
	deptByName map[string]uint
	prefByName map[string]uint
}

func (u *employeeUseCase) buildImportLookup(companyId uint) (*importLookup, error) {
	emps, err := u.repo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	empByCode := make(map[string]*model.Employee, len(emps))
	for i := range emps {
		empByCode[emps[i].StaffCode] = &emps[i]
	}

	depts, err := u.deptRepo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	deptByName := make(map[string]uint, len(depts))
	for _, d := range depts {
		deptByName[d.Name] = d.ID
	}

	prefs, err := u.prefRepo.GetAll()
	if err != nil {
		return nil, err
	}
	prefByName := make(map[string]uint, len(prefs))
	for _, p := range prefs {
		prefByName[p.Name] = p.ID
	}

	return &importLookup{empByCode: empByCode, deptByName: deptByName, prefByName: prefByName}, nil
}

func (u *employeeUseCase) applyRelated(repo repository.EmployeeRepository, companyId uint, empID uint, row service.EmployeeCSVRow, l *importLookup) error {
	if row.PostCode != "" || row.PrefectureName != "" || row.City != "" ||
		row.AddressLine1 != "" || row.AddressLine2 != "" || row.Tel != "" {
		var prefId *uint
		if row.PrefectureName != "" {
			if id, ok := l.prefByName[row.PrefectureName]; ok {
				prefId = &id
			}
		}
		addr := &model.EmployeeAddress{
			CompanyId:    companyId,
			EmployeeId:   empID,
			PostCode:     nilStr(row.PostCode),
			PrefectureId: prefId,
			City:         nilStr(row.City),
			AddressLine1: nilStr(row.AddressLine1),
			AddressLine2: nilStr(row.AddressLine2),
			Tel:          nilStr(row.Tel),
		}
		if err := repo.UpsertAddress(addr); err != nil {
			return err
		}
	}

	if len(row.Tenures) > 0 {
		tenures := make([]model.EmployeeTenures, 0, len(row.Tenures))
		for _, t := range row.Tenures {
			parsed, err := time.Parse("2006-01-02", t.JoinedOn)
			if err != nil {
				return err
			}
			tenure := model.EmployeeTenures{CompanyId: companyId, EmployeeId: empID, JoinedOn: parsed}
			if t.ResignationOn != "" {
				rt, err := time.Parse("2006-01-02", t.ResignationOn)
				if err != nil {
					return err
				}
				tenure.ResignationOn = &rt
			}
			if t.ResignationType != "" {
				tenure.ResignationType = &t.ResignationType
			}
			if t.Status != "" {
				tenure.Status = &t.Status
			}
			tenures = append(tenures, tenure)
		}
		if err := repo.ReplaceTenures(companyId, empID, tenures); err != nil {
			return err
		}
	}

	if len(row.DepartmentNames) > 0 {
		var deptIds []uint
		for _, name := range row.DepartmentNames {
			if id, ok := l.deptByName[name]; ok {
				deptIds = append(deptIds, id)
			}
		}
		if err := repo.UpdateDepartments(companyId, empID, deptIds); err != nil {
			return err
		}
	}
	return nil
}

func nilStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
