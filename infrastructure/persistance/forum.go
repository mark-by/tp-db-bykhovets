package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
	"strings"
)

type Forum struct {
	db *pgx.ConnPool
}

func newForum(db *pgx.ConnPool) *Forum {
	return &Forum{db}
}

func (f Forum) Create(forum *entity.Forum) error {
	tx, err := f.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO forums (slug, title, author) "+
		"VALUES ($1, $2, $3) "+
		"RETURNING title", forum.Slug, forum.Title, forum.Author)

	if err != nil {
		if strings.Contains(err.Error(), "forums_author_fkey") {
			EndTx(tx, err)
			return entityErrors.UserNotFound
		}
		logrus.Errorf("CREATE FORUM: %s", err.Error())
		return err
	}

	EndTx(tx, err)
	return nil
}

func (f Forum) GetBySlug(slug string) (*entity.Forum, error) {
	panic("implement me")
}

var _ repository.Forum = &Forum{}