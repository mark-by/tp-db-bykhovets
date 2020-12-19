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
	"strconv"
	"time"
)

type Post struct {
	db *pgx.ConnPool
}

func (p Post) Create(thread *entity.Thread, posts []entity.Post) error {
	if len(posts) == 0 {
		return nil
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { EndTx(tx, err) }()

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

	err = updateForumPostsCount(tx, thread.Forum, len(posts))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	tx, err = p.db.Begin()

	err = insertUsers(tx, thread.Forum, uniqAuthors(posts))
	if err != nil {
		logrus.Error(err)
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
		values += fmt.Sprintf("('%s', '%s'),", forum, author)
	}
	_, err := tx.Exec("INSERT INTO forums_users (forum, nickname) " +
		"VALUES " + values[:len(values)-1] + " ON CONFLICT DO NOTHING;")
	if err != nil {
		return err
	}
	return nil
}

func uniqAuthors(posts []entity.Post) map[string]bool {
	set := make(map[string]bool, len(posts))

	for _, post := range posts {
		set[post.Author] = true
	}
	return set
}

func (p Post) Get(id int64, related []string) (*entity.PostFull, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(tx, err) }()

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
	row := tx.QueryRow(selects+joins+"WHERE p.id = $1", id)

	err = row.Scan(vars...)

	if err != nil {
		if IsNotFoundErr(err) {
			return nil, entityErrors.PostNotFound
		}
		return nil, err
	}

	post.Created = created.Time.Format(time.RFC3339Nano)

	postFull := entity.PostFull{Post: &post}
	if authorRelated {
		if authorAbout.Valid {
			user.About = authorAbout.String
		}
		postFull.Author = &user
	}

	if threadRelated {
		if threadSlug.Valid {
			thread.Slug = threadSlug.String
		}
		thread.Created = threadCreated.Time.Format(time.RFC3339Nano)
		postFull.Thread = &thread
	}

	if forumRelated {
		postFull.Forum = &forum
	}

	return &postFull, nil
}

func (p Post) Update(post *entity.Post) error {
	if post.Message == "" {
		return entityErrors.NothingToUpdate
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() { EndTx(tx, err) }()

	created := sql.NullTime{}

	var prevMessage string
	err = tx.QueryRow("SELECT parent, message, created, author, thread, forum, is_edited FROM posts WHERE id = $1", post.ID).
		Scan(&post.Parent, &prevMessage, &created, &post.Author, &post.Thread, &post.Forum, &post.IsEdited)

	if err != nil {
		if IsNotFoundErr(err) {
			return entityErrors.PostNotFound
		}
		return err
	}

	if prevMessage == post.Message {
		post.Created = created.Time.Format(time.RFC3339Nano)
		return nil
	}

	err = tx.QueryRow("UPDATE posts SET message = $1, is_edited = true "+
		"WHERE id = $2 "+
		"RETURNING parent, created, author, thread, forum, is_edited", post.Message, post.ID).
		Scan(&post.Parent, &created, &post.Author, &post.Thread, &post.Forum, &post.IsEdited)

	if err != nil {
		if IsNotFoundErr(err) {
			return entityErrors.PostNotFound
		}
		return err
	}

	if created.Valid {
		post.Created = created.Time.Format(time.RFC3339Nano)
	}

	return nil
}

func (p Post) GetForThread(threadId int, desc bool, sortType string, since int, limit int) (entity.PostList, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(tx, err) }()

	selects := "SELECT p.id, p.message, p.is_edited, p.parent, p.created, p.author, p.thread, p.forum "

	limits := ""
	if limit != 0 && sortType != "parent_tree" {
		limits = fmt.Sprintf("LIMIT %d", limit)
	}

	from := "FROM posts AS p "

	sqlQuery := selects + from + getPostsWhere(threadId, sortType, desc, since, limit) +
		getPostsOrder(sortType, desc) + limits + ";"
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	posts := entity.PostList{}
	created := sql.NullTime{}
	for rows.Next() {
		post := entity.Post{}
		err := rows.Scan(&post.ID, &post.Message, &post.IsEdited, &post.Parent, &created, &post.Author, &post.Thread, &post.Forum)
		if err != nil {
			rows.Close()
			return nil, err
		}
		post.Created = created.Time.Format(time.RFC3339Nano)
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

func getPostsWhere(threadId int, sortType string, desc bool, since int, limit int) string {

	where := fmt.Sprintf("WHERE p.thread = %d ", threadId)
	symbol := ">"
	if desc {
		symbol = "<"
	}
	switch sortType {
	case "":
		fallthrough
	case "flat":
		if since == 0 {
			break
		}

		sinceAddition := fmt.Sprintf("AND p.id %s %d ", symbol, since)
		where += sinceAddition
	case "tree":
		if since == 0 {
			break
		}
		sinceAddition := fmt.Sprintf("AND p.path::bigint[] %s (SELECT path FROM posts WHERE id = %d)::bigint[] ",
			symbol, since)
		where += sinceAddition
	case "parent_tree":
		sinceAddition := ""
		if since != 0 {
			sinceAddition = fmt.Sprintf("AND p2.path[1] %s (SELECT path[1] FROM posts WHERE id = %d) ", symbol, since)
		}
		desStr := ""
		if desc {
			desStr = "DESC "
		}
		limits := ""
		if limit != 0 {
			limits = fmt.Sprintf("LIMIT %d", limit)
		}
		where = "JOIN " +
			"(SELECT * FROM posts AS p2 " +
			"WHERE p2.thread = " + strconv.Itoa(threadId) + " AND p2.parent = 0 " +
			sinceAddition +
			"ORDER BY p2.path " + desStr +
			limits +
			") AS prnt " +
			"ON prnt.path[1] = p.path[1] "
	default:
	}

	return where
}

func getPostsOrder(sortType string, desc bool) string {
	if sortType == "" {
		sortType = "flat"
	}
	descStr := ""
	if desc {
		descStr = "DESC "
	}
	switch sortType {
	case "flat":
		return fmt.Sprintf("ORDER BY p.created %s, p.id %s", descStr, descStr)
	case "tree":
		return fmt.Sprintf("ORDER BY p.path %s", descStr)
	case "parent_tree":
		return fmt.Sprintf("ORDER BY p.path[1] %s, p.path %s", descStr, descStr)
	default:
		return "ORDER BY p.created"
	}
}

func newPost(db *pgx.ConnPool) *Post {
	return &Post{db}
}

var _ repository.Post = &Post{}
