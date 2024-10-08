package domain

import (
	"bytes"
	"database/sql"
	"hello-cms/models"

	"github.com/go-playground/validator/v10"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type ContentDomain struct {
	db *sql.DB
}

func NewContentDomain(db *sql.DB) *ContentDomain {
	return &ContentDomain{db: db}
}

func (c *ContentDomain) GetContents() ([]models.Content, error) {
	query := `
		SELECT slug, title, posted_at 
		FROM entries
		WHERE is_visible=1
		ORDER BY posted_at DESC
		LIMIT 10;`

	contents := []models.Content{}
	rows, err := c.db.Query(query)
	if err != nil {
		return contents, err
	}

	defer rows.Close()

	for rows.Next() {
		var cn models.Content
		err := rows.Scan(&cn.Slug, &cn.Title, &cn.PostedAt)
		if err != nil {
			return contents, err
		}

		cn.Tags, err = cn.GetTags(c.db)
		if err != nil {
			return contents, err
		}
		contents = append(contents, cn)
	}

	return contents, nil
}

func (c *ContentDomain) GetTagedContents(tag string) ([]models.Content, error) {
	query := `
		SELECT slug, title, posted_at 
		FROM entries
		WHERE is_visible=1
			AND slug IN (SELECT slug FROM tags WHERE tag= ?)
		ORDER BY posted_at DESC
		;`

	contents := []models.Content{}
	rows, err := c.db.Query(query, tag)
	if err != nil {
		return contents, err
	}

	defer rows.Close()

	for rows.Next() {
		var cn models.Content
		err := rows.Scan(&cn.Slug, &cn.Title, &cn.PostedAt)
		if err != nil {
			return contents, err
		}

		cn.Tags, err = cn.GetTags(c.db)
		if err != nil {
			return contents, err
		}
		contents = append(contents, cn)
	}

	return contents, nil
}

func (c *ContentDomain) GetTags() ([]models.Tag, error) {
	tags := []models.Tag{}

	query := `
		SELECT tag, COUNT(0) AS _cnt 
		FROM tags 
		GROUP BY tag 
		ORDER BY _cnt DESC, tag ASC;`

	rows, err := c.db.Query(query)
	if err != nil {
		return tags, nil
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Tag
		err := rows.Scan(&t.Tag, &t.Count)
		if err != nil {
			return tags, err
		}

		tags = append(tags, t)
	}

	return tags, nil
}

func (c *ContentDomain) GetContent(slug string) (models.Content, *models.Response) {
	var content models.Content
	var response models.Response
	query := `SELECT 
	html, title, posted_at, slug
	FROM entries 
	WHERE slug=? AND is_visible=1`

	row := c.db.QueryRow(query, slug)
	err := row.Scan(&content.Html, &content.Title, &content.PostedAt, &content.Slug)
	if err != nil {
		response.StatusCode = 400
		response.Message = "Content Not Found."
		return content, &response
	}

	content.Tags, err = content.GetTags(c.db)
	if err != nil {
		response.StatusCode = 500
		response.Message = err.Error()
		return content, &response
	}

	return content, nil
}

func (c *ContentDomain) PostContent(body string) (err error) {

	gm := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	var body_buf bytes.Buffer
	context := parser.NewContext()
	err = gm.Convert([]byte(body), &body_buf, parser.WithContext(context))
	if err != nil {
		return err
	}

	html := body_buf.String()

	metaData := meta.Get(context)
	title := metaData["Title"]
	slug := metaData["Slug"]
	posted_at := metaData["Posted_at"]

	var tags []string
	tags_interface := metaData["Tags"]
	slice, ok := tags_interface.([]interface{})
	if ok {
		for _, tag := range slice {
			tags = append(tags, tag.(string))
		}
	}

	content := models.Content{
		Slug:     slug.(string),
		Title:    title.(string),
		Tags:     tags,
		PostedAt: posted_at.(string),
		Html:     html,
		Raw:      body,
	}

	var validate *validator.Validate = validator.New()

	validate.RegisterValidation("slug", models.ValidSlug)
	validate.RegisterValidation("posted_at", models.ValidPostedAt)

	err = validate.Struct(content)
	if err != nil {
		return err
	}

	err = content.Save(c.db)
	if err != nil {
		return err
	}

	return nil
}
