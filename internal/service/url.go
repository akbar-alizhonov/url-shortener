package service

import (
	"awesomeProject/internal/domain/url"
	"awesomeProject/internal/repositiries"
	"awesomeProject/pkg/logger"
	"context"
	"errors"
	"log/slog"
)

type UrlService interface {
	Save(ctx context.Context, urlToSave string, alias string) error
	List(ctx context.Context) ([]url.Url, error)
	Get(ctx context.Context, id int) (url.Url, error)
}

type urlService struct {
	repo      repositiries.UrlRepository
	generator AliasGenerator
	log       *slog.Logger
}

func NewUrlService(repo repositiries.UrlRepository, generator AliasGenerator, logger *slog.Logger) UrlService {
	return &urlService{repo: repo, generator: generator, log: logger}
}

func (s *urlService) Save(ctx context.Context, urlToSave string, alias string) error {
	if alias != "" {
		err := s.repo.Save(ctx, urlToSave, alias)
		if err != nil {
			s.log.Error(
				"failed to save url", slog.String("url", urlToSave),
				slog.String("alias", alias),
				slog.String("err", err.Error()),
				slog.String("request_id", logger.RequestIDFromContext(ctx)),
			)
			return err
		}
	}

	for i := 0; i < 5; i++ {
		alias = s.generator.Generate()
		err := s.repo.Save(ctx, urlToSave, alias)
		if err == nil {
			break
		}
		if errors.Is(err, url.ErrAliasTaken) {
			continue
		}
		s.log.Error(
			"failed to save url", urlToSave,
			slog.String("alias", alias),
			slog.String("err", err.Error()),
			slog.String("request_id", logger.RequestIDFromContext(ctx)),
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
