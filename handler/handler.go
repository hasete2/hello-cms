package handler

import (
	"hello-cms/domain"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	contentDomain domain.ContentDomain
	manageDomain  domain.ManageDomain
}

func NewHandler(c domain.ContentDomain, m domain.ManageDomain) *Handler {
	return &Handler{contentDomain: c, manageDomain: m}
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/init", h.init)
	e.GET("/contents", h.contents)
	e.GET("/contents/t/:tag", h.contents)
	e.GET("/tags", h.tags)
	e.GET("/c/:slug", h.content)
	e.POST("/content", h.post_content)
}
