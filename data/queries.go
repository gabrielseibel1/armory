package data

import (
	"database/sql"
)

type Queryer interface {
	Tables() ([]string, error)
}

type queries struct {
	db *sql.DB
}

func NewQueryer(db *sql.DB) Queryer {
	return &queries{db: db}
}

func (q *queries) Tables() ([]string, error) {
	rows, err := q.db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name;
	`)
	if err != nil {
		return []string{}, err
	}

	var tables []string
	var t string
	for rows.Next() {
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}