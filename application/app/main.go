package app

import (
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

func New(repositories *repository.Repositories) *application.App {
	return &application.App{
		Forum:   newForum(repositories),
		Post:    newPost(repositories),
		Service: newService(repositories),
		Thread:  newThread(repositories),
		User:    newUser(repositories),
	}
}
