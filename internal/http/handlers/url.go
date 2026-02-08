package handlers

import (
	"awesomeProject/internal/domain/url"
	"awesomeProject/internal/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type UrlHandler struct {
	serv service.UrlService
}

type urlSaveRequest struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

func NewUrlHandler(serv service.UrlService) *UrlHandler {
	return &UrlHandler{serv: serv}
}

func (h *UrlHandler) SaveUrl(c *echo.Context) error {
	var req urlSaveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.serv.SaveUrl(c.Request().Context(), req.Url, req.Alias)
	if err != nil {
		switch {
		case errors.Is(err, url.ErrAliasTaken):
			return c.JSON(
				http.StatusConflict,
				map[string]string{"error": "alias already taken"},
			)
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.NoContent(http.StatusCreated)
}
