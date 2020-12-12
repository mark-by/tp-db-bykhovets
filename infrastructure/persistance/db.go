package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
)

func New() *repository.Repositories {
	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{
			Database:             "forum",
			User:                 "forum",
			Password:             "123",
		},
		MaxConnections: 100,
	})

	if err != nil {
		logrus.Fatal("Fail to create repositories: %s", err)
	}

	return &repository.Repositories{
		Forum:   newForum(db),
		Post:    nil,
		Service: nil,
		Thread:  nil,
		User:    nil,
		Vote:    nil,
	}
}