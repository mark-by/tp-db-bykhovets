package application

import "github.com/mark-by/tp-db-bykhovets/domain/entity"

type App struct {
	Forum Forum
	Post Post
	Service Service
	Thread Thread
	User User
}

type Forum interface {
	Create(forum *entity.Forum) (*entity.Forum, error)
	Get(slug string) (*entity.Forum, error)
	CreateThread(slug string, thread *entity.Thread) (*entity.Thread, error)
	GetUsers(slug string, limit int32, since string, desc bool) ([]entity.User, error)
	GetThreads(slug string, limit int32, since string, desc bool) ([]entity.Thread, error)
}

type Post interface {
	Get(id int64, related []string) (*entity.PostFull, error)
	Update(id int64, post *entity.Post) (*entity.Post, error)
}

type Service interface {
	Clear() error
	Status() (*entity.Status, error)
}

type Thread interface {
	CreatePosts(slugOrId string, posts []entity.Post) ([]entity.Post, error)
	Get(slugOrId string) (*entity.Thread, error)
	Update(slugOrId string) (*entity.Thread, error)
	GetPosts(slugOrId string, limit int32, since int64, desc bool, sort string) ([]entity.Post, error)
	Vote(slugOrId string, vote *entity.Vote) (*entity.Thread, error)
}

type User interface {
	Create(user *entity.User) ([]entity.User, error)
	Get(nickname string) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
}