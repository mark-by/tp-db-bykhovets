package persistance

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
	"time"
)

type Post struct {
	db *pgx.ConnPool
}

func (p Post) Create(thread *entity.Thread, posts []entity.Post) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func(){EndTx(tx, err)}()

	values := ""
	for _, post := range posts {
		values += fmt.Sprintf("('%s', %d, '%s', %d, '%s', current_timestamp), ",
			post.Message,
			post.Parent,
			post.Author,
			thread.ID,
			thread.Forum)
	}
	sqlQuery := "INSERT INTO posts (message, parent, author, thread, forum, created) " +
		"VALUES " + values[:len(values)-2] +
		" RETURNING id, created;"

	rows, err := tx.Query(sqlQuery)

	if err != nil {
		return err
	}

	idx := 0
	created := sql.NullTime{}

	for rows.Next() {
		err = rows.Scan(&posts[idx].ID, &created)
		posts[idx].Created = created.Time.Format(time.RFC3339Nano)
		if err != nil {
			return err
		}
		posts[idx].Thread = thread.ID
		posts[idx].Forum = thread.Forum
		idx++
	}

	if idx == 0 {
		//ошибка в запросе
		_ = tx.Rollback()
		tx, err = p.db.Begin()
		if err != nil {
			return err
		}
		var id int64
		err = tx.QueryRow(sqlQuery).Scan(&id, &created)
		if err != nil {
			_ = tx.Rollback()
			switch true {
			case IsPostParentErr(err):
				return entityErrors.ParentNotExist
			case IsAuthorErr(err):
				return entityErrors.UserNotFound
			case IsForumErr(err):
				return entityErrors.ForumNotFound
			case IsThreadErr(err):
				return entityErrors.ThreadNotFound
			default:
				return err
			}
		}
	}

	err = insertUsers(tx, thread.Forum, uniqAuthors(posts))
	if err != nil {
		return err
	}

	err = updateForumPostsCount(tx, thread.Forum, len(posts))
	if err != nil {
		return err
	}

	return nil
}

func updateForumPostsCount(tx *pgx.Tx, forum string, num int) error {
	_, err := tx.Exec("UPDATE forums SET posts = (posts + $1) WHERE slug = $2;", num, forum)
	return err
}

func insertUsers(tx *pgx.Tx, forum string, authors map[string]bool) error {
	values := ""
	for author, _ := range authors {
		values += fmt.Sprintf("('%s', '%s'), ", forum, author)
	}
	_, err := tx.Exec("INSERT INTO forums_users (forum, nickname) " +
		"VALUES " + values[:len(values) - 2] + " ON CONFLICT DO NOTHING;")
	return err
}

func uniqAuthors(posts []entity.Post) map[string]bool {
	set := make(map[string]bool, len(posts))

	for _, post := range posts {
		set[post.Author] = true
	}
	return set
}

func (p Post) Get(id int64, related []string) (*entity.Post, *entity.User, *entity.Thread, *entity.Forum, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer func() {EndTx(tx, err)} ()

	post := entity.Post{ID: id}
	user := entity.User{}
	thread := entity.Thread{}
	forum := entity.Forum{}

	created := pgtype.Timestamptz{}
	selects := "SELECT p.parent, p.author, p.message, p.is_edited, p.forum, p.thread, p.created "
	joins := " FROM posts as p "

	authorRelated := false
	threadRelated := false
	forumRelated := false

	vars := []interface{}{
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&created,
	}

	threadCreated := pgtype.Timestamptz{}

	authorAbout := sql.NullString{}

	threadSlug := sql.NullString{}
	for _, rel := range related {
		switch rel {
		case "user":
			selects += ", c.nickname, c.fullname, c.about, c.email "
			joins += "JOIN customers AS c ON c.nickname = p.author "
			vars = append(vars, []interface{}{&user.NickName, &user.FullName, &authorAbout, &user.Email}...)
			authorRelated = true
		case "thread":
			selects += ", t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created "
			joins += "JOIN threads AS t ON t.id = p.thread "
			vars = append(vars, []interface{}{&thread.ID, &thread.Title,
				&thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &threadSlug, &threadCreated}...)
			threadRelated = true
		case "forum":
			selects += ", f.title, f.author, f.slug, f.posts, f.threads "
			joins += "JOIN forums AS f ON f.slug = p.forum "
			vars = append(vars, []interface{}{&forum.Title, &forum.Author, &forum.Slug, &forum.Posts, &forum.Threads}...)
			forumRelated = true
		}
	}
	logrus.Infof("SQL: %s WHERE p.id = %d", selects + joins, id)
	row := tx.QueryRow(selects + joins +
		"WHERE p.id = $1", id)

	err = row.Scan(vars...)

	if err != nil {
		if IsNotFoundErr(err) {
			return nil, nil, nil, nil, entityErrors.PostNotFound
		}
		return nil, nil, nil, nil, err
	}

	post.Created = created.Time.Format(time.RFC3339Nano)

	returnedUser := new(entity.User)
	if authorRelated {
		if authorAbout.Valid {
			returnedUser.About = authorAbout.String
		}
		returnedUser = &user
	} else {
		returnedUser = nil
	}

	returnedThread := new(entity.Thread)
	if threadRelated {
		if threadSlug.Valid {
			thread.Slug = threadSlug.String
		}
		thread.Created = threadCreated.Time.Format(time.RFC3339Nano)
		returnedThread = &thread
	} else {
		returnedThread = nil
	}

	returnedForum := new(entity.Forum)
	if forumRelated {
		returnedForum = &forum
	} else {
		returnedForum = nil
	}

	return &post, returnedUser, returnedThread, returnedForum, nil
}

func (p Post) Update(post *entity.Post) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {EndTx(tx, err)}()

	created := sql.NullTime{}

	err = tx.QueryRow("UPDATE posts SET message = $1, is_edited = true " +
		"WHERE id = $2 " +
		"RETURNING parent, created, author, thread, forum", post.Message, post.ID).
		Scan(&post.Parent, &created, &post.Author, &post.Thread, &post.Forum)

	if err != nil {
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (p Post) GetForThread(threadSlugOrId string, sortType string) ([]entity.Post, error) {
	panic("implement me")
}

func newPost(db *pgx.ConnPool) *Post {
	return &Post{db}
}

var _ repository.Post = &Post{}
