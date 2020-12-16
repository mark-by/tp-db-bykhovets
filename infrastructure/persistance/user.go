package persistance

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
	"strings"
)

type User struct {
	db *pgx.ConnPool
}

func (u User) Create(user *entity.User) (*entity.User, error) {
	about := &sql.NullString{}
	if user.About != "" {
		about.Valid = true
		about.String = user.About
	}
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func(){EndTx(tx, err)}()

	if _, err = tx.Exec("INSERT INTO customers (email, fullname, nickname, about) "+
		"VALUES ($1, $2, $3, $4)", user.Email, user.FullName, user.NickName, about); err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, entityErrors.UserAlreadyExist
		}
		return nil, err
	}
	return nil, nil
}

func (u User) GetForForum(slugForum string, since string, limit int, desc bool) ([]entity.User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {EndTx(tx, err)}()

	selects := fmt.Sprintf("SELECT u.nickname, u.fullname, u.email, u.about " +
		"FROM forums_users AS fs " +
		"JOIN customers as u ON fs.nickname = u.nickname " +
		"WHERE fs.forum = '%s' ", slugForum)

	symbol := ">"
	descStr := ""
	if desc {
		symbol = "<"
		descStr = "DESC "
	}

	sinceAddition := ""
	if since != "" {
		sinceAddition = fmt.Sprintf("AND fs.nickname %s '%s' ", symbol, since)
	}

	limits := ""
	if limit != 0 {
		limits = fmt.Sprintf("LIMIT %d ", limit)
	}

	order := fmt.Sprintf("ORDER BY u.nickname %s", descStr)
	sqlQuery := selects + sinceAddition + order + limits + ";"
	logrus.Info("SQL: ", sqlQuery)
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	var users []entity.User
	about := sql.NullString{}
	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(&user.NickName, &user.FullName, &user.Email, &about)
		if err != nil {
			return nil, err
		}
		if about.Valid {
			user.About = about.String
		}
		users = append(users, user)
	}
	return users, nil
}

func (u User) GetByNickName(nickname string) (*entity.User, error) {
	panic("implement me")
}

func (u User) Update(nickname string) (*entity.User, error) {
	panic("implement me")
}

func newUser(db *pgx.ConnPool) *User {
	return &User{db}
}

var _ repository.User = &User{}
