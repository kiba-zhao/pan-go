package repository

import (
	"database/sql"
	"iter"
)

func NewSeq2WithRows[T interface{}](rows *sql.Rows, db RepositoryDB) iter.Seq2[T, error] {

	return func(yield func(T, error) bool) {

		defer rows.Close()
		for rows.Next() {
			var model T
			err := db.ScanRows(rows, &model)
			if !yield(model, err) {
				break
			}
		}

	}

}
