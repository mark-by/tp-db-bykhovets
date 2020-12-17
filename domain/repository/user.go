package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type User interface {
	Create(user *entity.User) ([]entity.User, error)
	GetForForum(slugForum string, since string, limit int, desc bool) ([]entity.User, error)
	Get(user *entity.User) error
	Update(user *entity.User) error
}
