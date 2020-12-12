package persistance

import (
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

func EndTx(tx *pgx.Tx, err error) {
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Errorf("Fail to rollback: %s", err)
		}
		return
	}
	if err := tx.Commit(); err != nil {
		logrus.Errorf("Fail to commit: %s", err)
	}
}
