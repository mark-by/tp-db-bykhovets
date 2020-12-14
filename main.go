package main

import (
	"github.com/mark-by/tp-db-bykhovets/infrastructure/persistance"
	"github.com/sirupsen/logrus"
)

func main() {
	reps := persistance.New()

	post, user, thread, forum, err := reps.Post.Get(264, []string{"user", "forum", "thread"})

	if err != nil {
		logrus.Error("Error: ", err)
	}

	logrus.Infof("POST: %v", post)
	logrus.Infof("USER: %v", user)
	logrus.Infof("THREAD: %v", thread)
	logrus.Infof("FORUM: %v", forum)
	//handlers.ServeAPI(repositories)
}
