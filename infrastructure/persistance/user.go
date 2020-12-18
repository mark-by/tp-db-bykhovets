package persistance

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

type User struct {
	db *pgx.ConnPool
}

func (u User) Create(user *entity.User) ([]entity.User, error) {
	about := sql.NullString{}
	if user.About != "" {
		about.Valid = true
		about.String = user.About
	}
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(u.db, tx, err) }()

	_, err = tx.Exec("INSERT INTO customers (email, fullname, nickname, about) "+
		"VALUES ($1, $2, $3, $4)", user.Email, user.FullName, user.NickName, &about)

	if err == nil {
		return nil, nil
	}

	if !IsUniqErr(err) {
		return nil, err
	}

	err = tx.Rollback()
	if err != nil {
		return nil, err
	}
	tx, err = u.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT nickname, fullname, about, email "+
		"FROM customers "+
		"WHERE nickname = $1 or email = $2;", user.NickName, user.Email)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.NickName, &user.FullName, &about, &user.Email)
		if about.Valid {
			user.About = about.String
		}
		users = append(users, user)
	}

	return users, nil
}

func (u User) GetForForum(slugForum string, since string, limit int, desc bool) ([]entity.User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(u.db, tx, err) }()

	selects := fmt.Sprintf("SELECT u.nickname, u.fullname, u.email, u.about "+
		"FROM forums_users AS fs "+
		"JOIN customers as u ON fs.nickname = u.nickname "+
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
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := entity.UserList{}
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

func (u User) Get(user *entity.User) (err error) {
	tx, err := u.db.Begin()
	if err != nil {
		return
	}
	defer func() { EndTx(u.db, tx, err) }()

	about := sql.NullString{}
	err = tx.QueryRow("SELECT fullname, about, nickname, email FROM customers WHERE nickname = $1",
		user.NickName).Scan(&user.FullName, &about, &user.NickName, &user.Email)

	if err != nil {
		if IsNotFoundErr(err) {
			return entityErrors.UserNotFound
		}
		return
	}

	if about.Valid {
		user.About = about.String
	}

	return
}

func (u User) Update(user *entity.User) (err error) {
	tx, err := u.db.Begin()
	if err != nil {
		return
	}
	defer func() { EndTx(u.db, tx, err) }()

	columns := make([]string, 0, 3)
	values := make([]interface{}, 0, 3)
	if user.FullName != "" {
		columns = append(columns, "fullname")
		values = append(values, user.FullName)
	}
	if user.About != "" {
		columns = append(columns, "about")
		values = append(values, user.About)
	}
	if user.Email != "" {
		columns = append(columns, "email")
		values = append(values, user.Email)
	}

	if len(columns) == 0 {
		err = entityErrors.NothingToUpdate
		return
	}
	titles := updateTitles(columns)

	sqlRow := "UPDATE customers SET " + titles + fmt.Sprintf(" WHERE nickname = '%s' "+
		"RETURNING fullname, about, email", user.NickName)
	about := sql.NullString{}
	err = tx.QueryRow(sqlRow, values...).Scan(&user.FullName, &about, &user.Email)
	if err != nil {
		switch true {
		case IsNotFoundErr(err):
			return entityErrors.UserNotFound
		case IsUniqErr(err):
			return entityErrors.UserConflict
		default:
			return
		}
	}
	if about.Valid {
		user.About = about.String
	}
	return
}

func newUser(db *pgx.ConnPool) *User {
	return &User{db}
}

var _ repository.User = &User{}
