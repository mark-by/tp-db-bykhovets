package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

type Service struct {
	db *pgx.ConnPool
}

func (s Service) Clear() error {
	panic("implement me")
}

func (s Service) Info() (*entity.Status, error) {
	panic("implement me")
}



var _ repository.Service = &Service{}
