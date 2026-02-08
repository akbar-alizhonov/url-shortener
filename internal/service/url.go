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
	SaveUrl(ctx context.Context, urlToSave string, alias string) error
}

type urlService struct {
	repo      repositiries.UrlRepository
	generator AliasGenerator
	log       *slog.Logger
}

func NewUrlService(repo repositiries.UrlRepository, generator AliasGenerator, logger *slog.Logger) UrlService {
	return &urlService{repo: repo, generator: generator, log: logger}
}

func (s *urlService) SaveUrl(ctx context.Context, urlToSave string, alias string) error {
	if alias != "" {
		err := s.repo.Save(ctx, urlToSave, alias)
		if err != nil {
			s.log.Error(
				"failed to save url", urlToSave,
				"alias", alias,
				"err", err,
				"request_id", logger.RequestIDFromContext(ctx),
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
			"alias", alias,
			"err", err,
			"request_id", logger.RequestIDFromContext(ctx),
		)
		return err
	}

	return nil
}
