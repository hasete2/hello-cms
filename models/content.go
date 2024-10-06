package models

import (
	"database/sql"
	"regexp"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

type Content struct {
	Slug     string   `json:"slug" validate:"required,slug"`
	Title    string   `json:"title" validate:"required"`
	Tags     []string `json:"tags"`
	PostedAt string   `json:"posted_at" validate:"required,posted_at"`
	Html     string   `json:"html"`
	Raw      string   `json:"raw" validate:"required"`
}

func ValidSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	slugRegex := regexp.MustCompile(`^[a-zA-Z0-9]{1,}[a-zA-Z0-9-]{1,}[a-zA-Z0-9]{1}$`)
	return slugRegex.MatchString(slug)
}

func ValidPostedAt(fl validator.FieldLevel) bool {
	posted_at := fl.Field().String()
	postedatRegex := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	return postedatRegex.MatchString(posted_at)
}

func (c *Content) GetTags(db *sql.DB) (tags []string, err error) {
	query := "SELECT tag FROM tags  WHERE slug = ?"
	rows, err := db.Query(query, c.Slug)
	if err != nil {
		return tags, nil
	}
	defer rows.Close()

	var t string
	for rows.Next() {
		err := rows.Scan(&t)
		if err != nil {
			return tags, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

func (c *Content) Save(db *sql.DB) (err error) {
	query := "INSERT INTO entries(slug, raw, html, title, posted_at) VALUES(?, ?, ?, ?, ?);"
	_, err = db.Exec(query, c.Slug, c.Raw, c.Html, c.Title, c.PostedAt)
	if err != nil {
		return err
	}

	tag_query := "INSERT INTO tags(slug, tag) VALUES(?, ?);"
	for _, t := range c.Tags {
		_, err := db.Exec(tag_query, c.Slug, t)
		if err != nil {
			return err
		}
	}
	return err
}
