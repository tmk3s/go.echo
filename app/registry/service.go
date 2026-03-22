package registry

import (
	domainservice "app/domain/service"
	infraservice "app/infrastructure/service"
)

func (i *Registry) NewCsvService() domainservice.CsvService {
	return infraservice.NewCsvService()
}
