package persistance

import (
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
	"strings"
)

func EndTx(tx *pgx.Tx, err error) {
	if tx.Status() == -1 {
		return
	}
	if err != nil {
		_ = tx.Rollback()
		return
	}
	_ = tx.Commit()
}

func ConditionSlugOrId(slugOrId string) string {
	id, err := strconv.Atoi(slugOrId)
	condition := fmt.Sprintf("slug = '%s'", slugOrId)
	if err == nil {
		condition = fmt.Sprintf("id = %d", id)
	}
	return condition
}

func IsNotFoundErr(err error) bool {
	return strings.Contains(err.Error(), "no rows")
}

func IsPostParentErr(err error) bool {
	return strings.Contains(err.Error(), "Parent")
}

func IsAuthorErr(err error) bool {
	return strings.Contains(err.Error(), "author")
}

func IsForumErr(err error) bool {
	return strings.Contains(err.Error(), "forum")
}

func IsThreadErr(err error) bool {
	return strings.Contains(err.Error(), "thread")
}

func IsUniqErr(err error) bool {
	return strings.Contains(err.Error(), "unique")
}

func updateTitles(columns []string) string {
	values := ""
	for idx, column := range columns {
		values += fmt.Sprintf("%s = $%d,", column, idx + 1)
	}
	return values[:len(values) - 1]
}

type column struct {
	Title string
	Value interface{}
}