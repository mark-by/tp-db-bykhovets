package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Service interface {
	Clear() error
	Info() (*entity.Status, error)
}