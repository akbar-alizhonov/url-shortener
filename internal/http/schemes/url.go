package schemes

import (
	resp "awesomeProject/pkg/api/response"
)

type UrlBaseSchema struct {
	OriginalUrl string `json:"original_url"`
	Alias       string `json:"alias"`
}

type UrlGetSchema struct {
	Id int `json:"id"`
	UrlBaseSchema
	resp.Response
}

type UrlCreateSchema struct {
	UrlBaseSchema
}

type UrlUpdateSchema struct {
	Id     int    `json:"id"`
	NewUrl string `json:"new_url"`
	Alias  string `json:"alias"`
}
