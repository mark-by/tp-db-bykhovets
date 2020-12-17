package app

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

type Service struct {
	rep *repository.Repositories
}

func (s Service) Clear() error {
	return s.rep.Service.Clear()
}

func (s Service) Status() (*entity.Status, error) {
	return s.rep.Service.Info()
}

func newService(rep *repository.Repositories) *Service {
	return &Service{rep}
}

var _ application.Service = &Service{}