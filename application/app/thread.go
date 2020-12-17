package app

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"strconv"
)

type Thread struct {
	rep *repository.Repositories
}

func (t Thread) CreatePosts(slugOrId string, posts []entity.Post) error {
	thread, err := t.Get(slugOrId)

	if err != nil {
		return err
	}

	return t.rep.Post.Create(thread, posts)
}

func (t Thread) Get(slugOrId string) (*entity.Thread, error) {
	thread := entity.Thread{}
	num, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
	} else {
		thread.ID = int32(num)
	}
	err = t.rep.Thread.Get(&thread)
	if err != nil {
		return nil, err
	}
	return &thread, err
}

func (t Thread) Update(slugOrId string, thread *entity.Thread) error {
	num, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread.Slug = slugOrId
	} else {
		thread.ID = int32(num)
	}
	return t.rep.Thread.Update(thread)
}

func (t Thread) GetPosts(slugOrId string, limit int, since int, desc bool, sort string) ([]entity.Post, error) {
	thread, err := t.Get(slugOrId)
	if err != nil {
		return nil, err
	}
	return t.rep.Post.GetForThread(int(thread.ID),desc, sort, since, limit)
}

func (t Thread) Vote(slugOrId string, vote *entity.Vote) (*entity.Thread, error) {
	thread, err := t.Get(slugOrId)
	if err != nil {
		return nil, err
	}
	err = t.rep.Thread.Vote(vote, thread)
	if err != nil {
		return nil, err
	}
	return thread, err
}

func newThread(repositories *repository.Repositories) *Thread {
	return &Thread{rep: repositories}
}

var _ application.Thread = &Thread{}
