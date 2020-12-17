package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
)

type Service struct {
	db *pgx.ConnPool
}

func newService(db *pgx.ConnPool) *Service {
	return &Service{db}
}

func (s Service) Clear() (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer func() {EndTx(tx, err)}()

	_, err = tx.Exec("DELETE FROM customers")

	return
}


func (s Service) Info() (*entity.Status, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {EndTx(tx, err)}()
	rows, err := tx.Query(
		"SELECT count(*) FROM forums " +
			"UNION ALL " +
			"SELECT count(*) " +
			"FROM posts " +
			"UNION ALL " +
			"SELECT count(*) FROM threads " +
			"UNION ALL " +
			"SELECT count(*) FROM customers",
	)
	if err != nil {
		return nil, err
	}
	i := 0
	status := entity.Status{}
	for rows.Next() {
		var err error
		switch i {
		case 0:
			err = rows.Scan(&status.Forum)
		case 1:
			err = rows.Scan(&status.Post)
		case 2:
			err = rows.Scan(&status.Thread)
		case 3:
			err = rows.Scan(&status.User)
		}
		if err != nil {
			rows.Close()
			return nil, err
		}
		i++
	}
	rows.Close()
	return &status, err
}

var _ repository.Service = &Service{}
