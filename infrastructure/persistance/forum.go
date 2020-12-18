package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
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
	defer func() { EndTx(f.db, tx, err) }()

	err = tx.QueryRow("INSERT INTO forums (slug, title, author) "+
		"VALUES ($1, $2, (select nickname from customers where nickname = $3)) "+
		"RETURNING author", forum.Slug, forum.Title, forum.Author).Scan(&forum.Author)

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
	defer func() { EndTx(f.db, tx, err) }()

	err = tx.QueryRow("SELECT EXISTS (SELECT FROM forums WHERE slug = $1)", forumSlug).Scan(&exist)
	return
}

func (f Forum) GetBySlug(slug string) (*entity.Forum, error) {
	tx, err := f.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(f.db, tx, err) }()

	row := tx.QueryRow("SELECT f.posts, f.slug, f.threads, f.title, f.author "+
		"FROM forums as f "+
		"WHERE f.slug = $1", slug)

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
