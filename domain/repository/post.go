package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Post interface {
	Create(thread *entity.Thread, posts []entity.Post) error
	Get(id int64, related []string) (*entity.Post, *entity.User, *entity.Thread, *entity.Forum, error)
	Update(post *entity.Post) error
	GetForThread(threadSlugOrId string, sortType string) ([]entity.Post, error)
}
