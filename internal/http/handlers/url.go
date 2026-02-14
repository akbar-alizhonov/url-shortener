package handlers

import (
	"awesomeProject/internal/domain/url"
	"awesomeProject/internal/http/schemes"
	"awesomeProject/internal/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type UrlHandler struct {
	serv service.UrlService
}

func NewUrlHandler(serv service.UrlService) *UrlHandler {
	return &UrlHandler{serv: serv}
}

func (h *UrlHandler) SaveUrl(c *echo.Context) error {
	var req schemes.UrlCreateSchema
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.serv.Save(c.Request().Context(), req.OriginalUrl, req.Alias)
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

func (h *UrlHandler) ListUrls(c *echo.Context) error {
	urls, err := h.serv.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := make([]schemes.UrlGetSchema, len(urls))
	for idx, u := range urls {
		resp[idx] = schemes.UrlGetSchema{
			Id: u.Id,
			UrlBaseSchema: schemes.UrlBaseSchema{
				OriginalUrl: u.OriginalUrl,
				Alias:       u.Alias,
			},
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UrlHandler) Redirect(c *echo.Context) error {
	id, err := echo.PathParam[int](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	u, err := h.serv.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.Redirect(http.StatusFound, u.OriginalUrl)
}

func (h *UrlHandler) Update(c *echo.Context) error {
	var req schemes.UrlUpdateSchema
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.serv.Update(c.Request().Context(), req.Id, req.NewUrl, req.Alias)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
