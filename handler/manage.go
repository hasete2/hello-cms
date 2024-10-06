package handler

import (
	"hello-cms/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) init(c echo.Context) error {

	err := h.manageDomain.Init()
	if err != nil {
		response := models.Response{
			StatusCode: 400,
			Message:    err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	response := models.Response{
		StatusCode: 200,
		Message:    "OK",
	}
	return c.JSON(http.StatusOK, response)
}
