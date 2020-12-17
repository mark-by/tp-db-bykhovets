package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Post interface {
	Create(thread *entity.Thread, posts []entity.Post) error
	Get(id int64, related []string) (*entity.PostFull, error)
	Update(post *entity.Post) error
	GetForThread(id int, desc bool, sortType string, since int, limit int) (entity.PostList, error)
}
