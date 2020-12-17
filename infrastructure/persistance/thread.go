package persistance

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"time"
)

type Thread struct {
	db *pgx.ConnPool
}

func (t Thread) Create(thread *entity.Thread) error {
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
		"RETURNING id;", slug, thread.Title, thread.Message, thread.Author, thread.Forum, created).Scan(&id)
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
	return nil
}

func (t Thread) GetForForum(forumSlug string, since string, limit int, desc bool) ([]entity.Thread, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {EndTx(tx, err)}()

	descString := ""
	symbol := ">="
	if desc {
		descString = "DESC "
		symbol = "<="
	}

	selects := "SELECT t.id, t.slug, t.message, t.created, t.votes, t.author, t.forum " +
		"FROM threads AS t "
	where := fmt.Sprintf("WHERE t.forum = '%s' ", forumSlug)
	order := fmt.Sprintf("ORDER BY t.created %s", descString)
	limits := ""
	if limit != 0 {
		limits = fmt.Sprintf("LIMIT %d", limit)
	}

	if since != "" {
		where += fmt.Sprintf("AND t.created %s '%s' ", symbol, since)
	}

	sqlQuery := selects + where + order + limits
	rows, err := tx.Query(sqlQuery)

	if err != nil {
		return nil, err
	}
	var threads []entity.Thread

	slug := sql.NullString{}
	created := sql.NullTime{}
	for rows.Next() {
		thread := entity.Thread{}
		err := rows.Scan(&thread.ID, &slug, &thread.Message,
			&created, &thread.Votes, &thread.Author, &thread.Forum)
		if err != nil {
			return nil, err
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		thread.Created = created.Time.Format(time.RFC3339Nano)
		threads = append(threads, thread)
	}
	return threads, nil
}

func (t Thread) Get(thread *entity.Thread) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer func() {EndTx(tx, err)}()

	selects := "SELECT id, title, author, forum, message, votes, slug, created FROM threads "
	where := fmt.Sprintf("WHERE id = %d;", thread.ID)
	if thread.Slug != "" {
		where = fmt.Sprintf("WHERE slug = '%s';", thread.Slug)
	}

	slug := sql.NullString{}
	created := sql.NullTime{}
	err = tx.QueryRow(selects + where).Scan(&thread.ID, &thread.Title, &thread.Author,
		&thread.Forum, &thread.Message, &thread.Votes, &slug, &created)
	if err != nil {
		return err
	}
	if slug.Valid {
		thread.Slug = slug.String
	}
	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}
	return nil
}

func (t Thread) Update(thread *entity.Thread) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer func() {EndTx(tx, err)}()
	updateColumns := make([]string, 0, 2)
	values := make([]interface{}, 0, 2)
	if thread.Title != "" {
		updateColumns = append(updateColumns, "title")
		values = append(values, thread.Title)
	}
	if thread.Message != "" {
		updateColumns = append(updateColumns, "message")
		values = append(values, thread.Message)
	}
	titles := updateTitles(updateColumns)
	where := fmt.Sprintf(" WHERE id = %d ", thread.ID)
	if thread.Slug != "" {
		where = fmt.Sprintf(" WHERE slug = '%s' ", thread.Slug)
	}
	created := sql.NullTime{}
	slug := sql.NullString{}

	sqlRow := "UPDATE threads SET " + titles + where +
		"RETURNING id, title, message, created, slug, author, forum, votes;"
	err = tx.QueryRow(sqlRow, values...).Scan(&thread.ID, &thread.Title,
		&thread.Message, &created, &slug, &thread.Author, &thread.Forum, &thread.Votes)
	if err != nil {
		if IsNotFoundErr(err) {
			return entityErrors.ThreadNotFound
		}
		return err
	}
	if created.Valid {
		thread.Created = created.Time.Format(time.RFC3339Nano)
	}
	if slug.Valid {
		thread.Slug = slug.String
	}
	return nil
}

func (t Thread) Vote(vote *entity.Vote, thread *entity.Thread) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	defer func() {EndTx(tx, err)}()
	var voice int32
	err = tx.QueryRow("SELECT voice FROM votes AS v WHERE v.thread = $1", thread.ID).Scan(&voice)
	if err == nil && voice == vote.Voice {
		return nil
	}

	_, err = tx.Exec("INSERT INTO VOTES (voice, author, thread) VALUES ($1, $2, $3) " +
		"ON CONFLICT (author, thread) DO UPDATE SET voice = EXCLUDED.voice;", vote.Voice, vote.Author, thread.ID)

	if err != nil {
		switch true {
		case IsAuthorErr(err):
			return entityErrors.UserNotFound
		case IsThreadErr(err):
			return entityErrors.ThreadNotFound
		default:
			return err
		}
	}

	thread.Votes += vote.Voice - voice

	return nil
}

func newThread(db *pgx.ConnPool) *Thread {
	return &Thread{db}
}

var _ repository.Thread = &Thread{}
