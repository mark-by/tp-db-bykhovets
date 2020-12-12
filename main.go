package main

import (
	"github.com/mark-by/tp-db-bykhovets/application/persistanceApp"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/infrastructure/persistance"
	"github.com/sirupsen/logrus"
)

func main() {
	reps := persistance.New()
	app := persistanceApp.New(reps)

	err := app.Forum.Create(&entity.Forum{
		Slug:    "freeky",
		Author:  "biem",
		Title:   "Fuck you",
	})
	if err != nil {
		logrus.Error(err)
	}

	//handlers.ServeAPI(repositories)
}