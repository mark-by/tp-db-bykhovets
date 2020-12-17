package app

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
)

type Forum struct {
	rep *repository.Repositories
}

func newForum(repositories *repository.Repositories) *Forum {
	return &Forum{rep: repositories}
}

func (f Forum) Create(forum *entity.Forum) error {
	err := f.rep.Forum.Create(forum)
	if err == entityErrors.ForumAlreadyExist {
		existedForum, err := f.rep.Forum.GetBySlug(forum.Slug)
		if err != nil {
			logrus.WithField("action", "create_forum").Error(err)
			return err
		}
		*forum = *existedForum
		return entityErrors.ForumAlreadyExist
	}
	return err
}

func (f Forum) Get(slug string) (*entity.Forum, error) {
	return f.rep.Forum.GetBySlug(slug)
}

func (f Forum) CreateThread(slug string, thread *entity.Thread) error {
	thread.Forum = slug
	err := f.rep.Thread.Create(thread)
	if err == entityErrors.ThreadAlreadyExist {
		err = f.rep.Thread.Get(thread)
		if err != nil {
			return err
		}
		return entityErrors.ThreadAlreadyExist
	}
	return nil
}

func (f Forum) GetUsers(slug string, limit int, since string, desc bool) (entity.UserList, error) {
	ok, _ := f.rep.Forum.Exists(slug)
	if !ok {
		return nil, entityErrors.ForumNotFound
	}
	return f.rep.User.GetForForum(slug, since, limit, desc)
}

func (f Forum) GetThreads(slug string, limit int, since string, desc bool) (entity.ThreadList, error) {
	ok, _ := f.rep.Forum.Exists(slug)
	if !ok {
		return nil, entityErrors.ForumNotFound
	}
	return f.rep.Thread.GetForForum(slug, since, limit, desc)
}

var _ application.Forum = &Forum{}
