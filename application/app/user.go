package app

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

type User struct {
	rep *repository.Repositories
}

func (u User) Create(user *entity.User) (entity.UserList, error) {
	return u.rep.User.Create(user)
}

func (u User) Get(nickname string) (*entity.User, error) {
	user := entity.User{NickName: nickname}
	err := u.rep.User.Get(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u User) Update(user *entity.User) error {
	return u.rep.User.Update(user)
}

func newUser(rep *repository.Repositories) *User {
	return &User{rep}
}

var _ application.User = &User{}
