package entityErrors

import "github.com/pkg/errors"

var (
	UserNotFound     = errors.New("user not found")
	UserAlreadyExist = errors.New("user already exist")
	UserConflict     = errors.New("user conflict")

	NothingToUpdate = errors.New("nothing to update")

	ForumNotFound     = errors.New("forum not exist")
	ForumAlreadyExist = errors.New("forum already exist")

	ThreadNotFound     = errors.New("thread not found")
	ThreadAlreadyExist = errors.New("thread already exist")

	PostNotFound   = errors.New("post not found")
	ParentNotExist = errors.New("parent post not exist in thread")
)
