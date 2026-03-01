package service

import (
	"awesomeProject/internal/domain/url"
	"awesomeProject/internal/repositiries"
	"awesomeProject/pkg/logger"
	"context"
	"errors"
	"log/slog"
	"strings"
)

type UrlService interface {
	Save(ctx context.Context, urlToSave, alias string) error
	List(ctx context.Context) ([]url.Url, error)
	Get(ctx context.Context, id int) (url.Url, error)
	Update(ctx context.Context, id int, newUrl, alias string) error
	Delete(ctx context.Context, id int) error
}

type urlService struct {
	repo      repositiries.UrlRepository
	generator AliasGenerator
	log       *slog.Logger
	baseUrl   string
}

func NewUrlService(
	repo repositiries.UrlRepository,
	generator AliasGenerator,
	logger *slog.Logger,
	baseUrl string,
) UrlService {
	return &urlService{
		repo:      repo,
		generator: generator,
		log:       logger,
		baseUrl:   baseUrl,
	}
}

func (s *urlService) Save(ctx context.Context, urlToSave, alias string) error {
	log := s.log.With(
		slog.String("url", urlToSave),
		slog.String("alias", alias),
		slog.String("request_id", logger.RequestIDFromContext(ctx)),
	)

	if alias != "" {
		shortUrl := s.BuildShortUrl(s.baseUrl, alias)
		err := s.repo.Save(ctx, urlToSave, shortUrl)
		if err != nil {
			log.Error(
				"failed to save url",
				slog.String("err", err.Error()),
			)
			return err
		}
		return nil
	}

	for i := 0; i < 5; i++ {
		alias = s.generator.Generate()
		shortUrl := s.BuildShortUrl(s.baseUrl, alias)
		err := s.repo.Save(ctx, urlToSave, shortUrl)
		if err == nil {
			break
		}
		if errors.Is(err, url.ErrAliasTaken) {
			continue
		}
		log.Error(
			"failed to save url",
			slog.String("err", err.Error()),
		)
		return err
	}

	return nil
}

func (s *urlService) List(ctx context.Context) ([]url.Url, error) {
	urls, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *urlService) Get(ctx context.Context, id int) (url.Url, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return url.Url{}, err
	}

	return u, nil
}

func (s *urlService) BuildShortUrl(baseUrl, code string) string {
	baseUrl = strings.TrimRight(baseUrl, "/")
	return baseUrl + "/" + code
}

func (s *urlService) Update(ctx context.Context, id int, newUrl, alias string) error {
	var shortUrl string
	if alias != "" {
		shortUrl = s.BuildShortUrl(s.baseUrl, alias)
	}

	err := s.repo.Update(ctx, id, newUrl, shortUrl)
	if err != nil {
		s.log.Error(
			"failed to update url", slog.String("url", newUrl),
			slog.String("alias", alias),
		)
		return err
	}

	return nil
}

func (s *urlService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
