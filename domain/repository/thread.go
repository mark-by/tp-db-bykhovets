package repository

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type Thread interface {
	Create(thread *entity.Thread) error
	GetForForum(forumSlug string, since string, limit int, desc bool) ([]entity.Thread, error)
	Get(thread *entity.Thread) error
	Update(thread *entity.Thread) error
	Vote(vote *entity.Vote, thread *entity.Thread) error
}
