package repositiries

import (
	"awesomeProject/internal/domain/url"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRepository interface {
	Save(ctx context.Context, urlToSave string, alias string) error
	List(ctx context.Context) ([]url.Url, error)
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

func (r *urlRepository) List(ctx context.Context) ([]url.Url, error) {
	sql, args, err := sq.
		Select("id", "original_url", "alias").From("url").
		OrderBy("created_at").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []url.Url
	for rows.Next() {
		var u url.Url
		if err := rows.Scan(&u.Id, &u.OriginalUrl, &u.Alias); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
