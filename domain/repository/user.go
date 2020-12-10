package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type User interface {
	Create(user *entity.User) (*entity.User, error)
	GetForForum(slugForum string) ([]*entity.User, error)
	GetByNickName(nickname string) (*entity.User, error)
	Update(nickname string) (*entity.User, error)
}