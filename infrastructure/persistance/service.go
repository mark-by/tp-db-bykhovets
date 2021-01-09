package persistance

import (
	"github.com/jackc/pgx"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/sirupsen/logrus"
)

type Service struct {
	db *pgx.ConnPool
}

func newService(db *pgx.ConnPool) *Service {
	service := Service{db}
	err := service.Prepare()
	if err != nil {
		logrus.Error(err.Error())
	}
	return &service
}

func (s Service) Clear() (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer func() { EndTx(tx, err) }()

	_, err = tx.Exec("DELETE FROM customers")

	return
}

func (s Service) Prepare() error {
	_, err := s.db.Prepare("info", "SELECT count(*) FROM forums "+
		"UNION ALL "+
		"SELECT count(*) "+
		"FROM posts "+
		"UNION ALL "+
		"SELECT count(*) FROM threads "+
		"UNION ALL "+
		"SELECT count(*) FROM customers")

	if err != nil {
		return err
	}

	return nil
}

func (s Service) Info() (*entity.Status, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { EndTx(tx, err) }()
	rows, err := tx.Query("info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
			return nil, err
		}
		i++
	}
	return &status, err
}

var _ repository.Service = &Service{}
