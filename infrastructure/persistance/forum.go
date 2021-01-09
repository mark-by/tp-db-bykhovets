package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
)

type Forum struct {
	db *pgx.ConnPool
}

func newForum(db *pgx.ConnPool) *Forum {
	forum := Forum{db}
	err := forum.Prepare()
	if err != nil {
		logrus.Error(err.Error())
	}
	return &forum
}

func (f Forum) Prepare() error {
	if _, err := f.db.Prepare("createForum", "INSERT INTO forums (slug, title, author) "+
		"VALUES ($1, $2, (select nickname from customers where nickname = $3)) "+
		"RETURNING author"); err != nil {
		return err
	}

	if _, err := f.db.Prepare("existsForum", "SELECT EXISTS (SELECT FROM forums WHERE slug = $1)"); err != nil {
		return err
	}

	if _, err := f.db.Prepare("getForumBySlug", "SELECT f.posts, f.slug, f.threads, f.title, f.author "+
		"FROM forums as f WHERE f.slug = $1"); err != nil {
		return err
	}

	return nil
}

func (f Forum) Create(forum *entity.Forum) error {
	tx, err := f.db.Begin()
	if err != nil {
		return err
	}
	defer func() { EndTx(tx, err) }()

	err = tx.QueryRow("createForum", forum.Slug, forum.Title, forum.Author).Scan(&forum.Author)

	if err != nil {
		switch true {
		case IsAuthorErr(err):
			return entityErrors.UserNotFound
		case IsUniqErr(err):
			return entityErrors.ForumAlreadyExist
		default:
			return err
		}
	}

	return nil
}

func (f Forum) Exists(forumSlug string) (exist bool, err error) {
	tx, err := f.db.Begin()
	if err != nil {
		return
	}
	defer func() { EndTx(tx, err) }()

	err = tx.QueryRow("existsForum", forumSlug).Scan(&exist)
	return
}

func (f Forum) GetBySlug(slug string) (*entity.Forum, error) {
	tx, err := f.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(tx, err) }()

	row := tx.QueryRow("getForumBySlug", slug)

	forum := entity.Forum{
		Slug: slug,
	}
	if err := row.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.Author); err != nil {
		if IsNotFoundErr(err) {
			return nil, entityErrors.ForumNotFound
		}
		return nil, err
	}
	return &forum, nil
}

var _ repository.Forum = &Forum{}
