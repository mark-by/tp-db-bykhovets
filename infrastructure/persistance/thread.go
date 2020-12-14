package persistance

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"time"
)

type Thread struct {
	db *pgx.ConnPool
}

func (t Thread) Create(forumSlug string, thread *entity.Thread) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer func() {EndTx(tx, err)} ()

	created := sql.NullTime{Time: time.Now(), Valid: true}
	if thread.Created != "" {
		if created.Time, err = time.Parse(time.RFC3339Nano, thread.Created); err != nil {
			return err
		}
	}

	slug := sql.NullString{}
	if thread.Slug != "" {
		slug.String = thread.Slug
		slug.Valid = true
	}

	var id int32
	err = tx.QueryRow("INSERT INTO threads (slug, title, message, author, forum, created) "+
		"VALUES ($1, $2, $3, $4, $5, $6) "+
		"RETURNING id;", slug, thread.Title, thread.Message, thread.Author, forumSlug, created).Scan(&id)
	if err != nil {
		switch true {
		case IsUniqErr(err):
			return entityErrors.ThreadAlreadyExist
		case IsAuthorErr(err):
			return entityErrors.UserNotFound
		case IsForumErr(err):
			return entityErrors.ForumNotFound
		default:
			return err
		}
	}

	thread.ID = id
	if thread.Created == "" {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}
	thread.Forum = forumSlug
	return nil
}

func (t Thread) GetForForum(forumSlug string) ([]entity.Thread, error) {
	panic("implement me")
}

func (t Thread) GetByID(forumID int32) (*entity.Thread, error) {
	panic("implement me")
}

func (t Thread) GetBySlug(forumSlug string) (*entity.Thread, error) {
	panic("implement me")
}

func (t Thread) Update(thread *entity.Thread) error {
	panic("implement me")
}

func (t Thread) Voice(thread *entity.Thread) error {
	panic("implement me")
}

func newThread(db *pgx.ConnPool) *Thread {
	return &Thread{db}
}

var _ repository.Thread = &Thread{}
