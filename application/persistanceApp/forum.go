package persistanceApp

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
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
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (f Forum) Get(slug string) (*entity.Forum, error) {
	panic("implement me")
}

func (f Forum) CreateThread(slug string, thread *entity.Thread) (*entity.Thread, error) {
	panic("implement me")
}

func (f Forum) GetUsers(slug string, limit int32, since string, desc bool) ([]entity.User, error) {
	panic("implement me")
}

func (f Forum) GetThreads(slug string, limit int32, since string, desc bool) ([]entity.Thread, error) {
	panic("implement me")
}

var _ application.Forum = &Forum{}
