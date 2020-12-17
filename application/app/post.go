package app

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

type Post struct {
	rep *repository.Repositories
}

func (p Post) Get(id int64, related []string) (*entity.PostFull, error) {
	return p.rep.Post.Get(id, related)
}

func (p Post) Update(post *entity.Post) error {
	return p.rep.Post.Update(post)
}

func newPost(repositories *repository.Repositories) *Post {
	return &Post{rep: repositories}
}

var _ application.Post = &Post{}

