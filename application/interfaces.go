package application

import (
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
)

type App struct {
	Forum   Forum
	Post    Post
	Service Service
	Thread  Thread
	User    User
}

type Forum interface {
	Create(forum *entity.Forum) error
	Get(slug string) (*entity.Forum, error)
	CreateThread(slug string, thread *entity.Thread) error
	GetUsers(slug string, limit int, since string, desc bool) (entity.UserList, error)
	GetThreads(slug string, limit int, since string, desc bool) (entity.ThreadList, error)
}

type Post interface {
	Get(id int64, related []string) (*entity.PostFull, error)
	Update(post *entity.Post) error
}

type Service interface {
	Clear() error
	Status() (*entity.Status, error)
}

type Thread interface {
	CreatePosts(slugOrId string, posts []entity.Post) error
	Get(slugOrId string) (*entity.Thread, error)
	Update(slugOrId string, thread *entity.Thread) error
	GetPosts(slugOrId string, limit int, since int, desc bool, sort string) (entity.PostList, error)
	Vote(slugOrId string, vote *entity.Vote) (*entity.Thread, error)
}

type User interface {
	Create(user *entity.User) (entity.UserList, error)
	Get(nickname string) (*entity.User, error)
	Update(user *entity.User) error
}
