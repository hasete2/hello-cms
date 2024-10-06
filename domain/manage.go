package domain

import (
	"database/sql"
)

type ManageDomain struct {
	db *sql.DB
}

func NewManageDomain(db *sql.DB) *ManageDomain {
	return &ManageDomain{db: db}
}

func (m *ManageDomain) Init() (err error) {
	drop_query := `
		DROP TABLE IF EXISTS "entries";
		DROP TABLE IF EXISTS "tags";`
	_, err = m.db.Exec(drop_query)
	if err != nil {
		return err
	}

	create_query := `
		CREATE TABLE IF NOT EXISTS "entries" (
			"slug"	TEXT NOT NULL,
			"raw"	TEXT NOT NULL,
			"html"	TEXT NOT NULL,
			"title"	TEXT NOT NULL,
			"posted_at"	TEXT NOT NULL,
			"is_visible"	TEXT NOT NULL DEFAULT 1,
			PRIMARY KEY("slug")
		);
		CREATE TABLE IF NOT EXISTS "tags" (
			"slug"	TEXT NOT NULL,
			"tag"	TEXT NOT NULL
		);`
	_, err = m.db.Exec(create_query)
	if err != nil {
		return err
	}

	return nil
}

func (m *ManageDomain) Drop() (err error) {

	return nil
}
