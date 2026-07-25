package usecase

func (u *employeeUseCase) ExportCSV(companyId uint) ([]byte, error) {
	employees, err := u.repo.GetListForExport(companyId)
	if err != nil {
		return nil, err
	}
	return u.csvService.GenerateEmployeeCSV(employees)
}
