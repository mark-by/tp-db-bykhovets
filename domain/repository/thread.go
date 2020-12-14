package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Thread interface {
	Create(forumSlug string, thread *entity.Thread) error
	GetForForum(forumSlug string) ([]entity.Thread, error)
	GetByID(forumID int32) (*entity.Thread, error)
	GetBySlug(forumSlug string) (*entity.Thread, error)
	Update(thread *entity.Thread) error
	Voice(thread *entity.Thread) error
}
