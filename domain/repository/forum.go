package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Forum interface {
	Create(forum *entity.Forum) error
	GetBySlug(slug string) (*entity.Forum, error)
	Exists(forumSlug string) (exist bool, err error)
}
