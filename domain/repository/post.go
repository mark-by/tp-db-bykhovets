package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Post interface {
	CreateBySlug(threadSlug string, posts []*entity.Post) ([]*entity.Post, error)
	CreateByID(threadID int32, posts []*entity.Post) ([]*entity.Post, error)
	Get(id int64) (*entity.PostFull, error)
	Update(post *entity.Post) (*entity.Post, error)
	GetForThreadBySlug(threadSlug string, sortType string) ([]*entity.Post, error)
	GetForThreadByID(threadID int32, sortType string) ([]*entity.Post, error)
}
