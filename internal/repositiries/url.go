package repositiries

import (
	"awesomeProject/internal/domain/url"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRepository interface {
	Save(ctx context.Context, urlToSave string, alias string) error
}

type urlRepository struct {
	pool *pgxpool.Pool
}

func NewUrlRepository(pool *pgxpool.Pool) UrlRepository {
	return &urlRepository{pool: pool}
}

func (r *urlRepository) Save(ctx context.Context, urlToSave string, alias string) error {
	sql, args, err := sq.
		Insert("url").Columns("original_url", "alias").
		Values(urlToSave, alias).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, sql, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return url.ErrAliasTaken
		}
		return err
	}

	return nil
}
