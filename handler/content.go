package handler

import (
	"hello-cms/models"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) contents(c echo.Context) error {
	var err error
	var contents []models.Content

	tag := c.Param("tag")
	if len(tag) > 0 {
		contents, err = h.contentDomain.GetTagedContents(tag)
	} else {
		contents, err = h.contentDomain.GetContents()
	}

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, contents)
}

func (h *Handler) tags(c echo.Context) error {
	tags, err := h.contentDomain.GetTags()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tags)
}

func (h *Handler) content(c echo.Context) error {
	slug := c.Param("slug")
	content, err := h.contentDomain.GetContent(slug)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, content)
}

func (h *Handler) post_content(c echo.Context) error {

	var buf []byte
	var body string
	buf, _ = io.ReadAll(c.Request().Body)
	body = string(buf)

	err := h.contentDomain.PostContent(body)
	if err != nil {
		response := models.Response{
			StatusCode: 400,
			Message:    err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	return c.String(http.StatusOK, "Hello, World!")
}
