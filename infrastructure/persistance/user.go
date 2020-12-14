package persistance

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
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

func (u User) GetForForum(slugForum string) ([]*entity.User, error) {
	panic("implement me")
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
