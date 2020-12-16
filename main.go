package main

import (
	"github.com/mark-by/tp-db-bykhovets/infrastructure/persistance"
	"github.com/sirupsen/logrus"
)

func main() {
	reps := persistance.New()

	forum := "forum2"
	ok, err := reps.Forum.Exists(forum)

	if err != nil {
		logrus.Error(err)
	}

	if !ok {
		logrus.Info("FORUM DOES NOT EXIST")
		return
	}

	users, err := reps.User.GetForForum(forum, "", 1, false)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("Users: %v", users)

	//handlers.ServeAPI(repositories)
}
