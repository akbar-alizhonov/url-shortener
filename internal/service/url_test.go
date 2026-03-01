package service

import (
	"awesomeProject/internal/domain/url"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
)

// --- Mocks ---

type mockRepo struct {
	saveFn   func(ctx context.Context, urlToSave, alias string) error
	listFn   func(ctx context.Context) ([]url.Url, error)
	getFn    func(ctx context.Context, id int) (url.Url, error)
	updateFn func(ctx context.Context, id int, newUrl, alias string) error
	deleteFn func(ctx context.Context, id int) error
}

func (m *mockRepo) Save(ctx context.Context, urlToSave, alias string) error {
	return m.saveFn(ctx, urlToSave, alias)
}

func (m *mockRepo) List(ctx context.Context) ([]url.Url, error) {
	return m.listFn(ctx)
}

func (m *mockRepo) Get(ctx context.Context, id int) (url.Url, error) {
	return m.getFn(ctx, id)
}

func (m *mockRepo) Update(ctx context.Context, id int, newUrl, alias string) error {
	return m.updateFn(ctx, id, newUrl, alias)
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

type mockGenerator struct {
	aliases []string
	index   int
}

func (m *mockGenerator) Generate() string {
	alias := m.aliases[m.index]
	m.index++
	return alias
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// --- Save tests ---

func TestSave_WithAlias_Success(t *testing.T) {
	repo := &mockRepo{
		saveFn: func(ctx context.Context, urlToSave, alias string) error {
			if urlToSave != "https://example.com" {
				t.Errorf("unexpected url: %s", urlToSave)
			}
			if alias != "http://localhost/my-alias" {
				t.Errorf("unexpected alias: %s", alias)
			}
			return nil
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Save(context.Background(), "https://example.com", "my-alias")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestSave_WithAlias_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mockRepo{
		saveFn: func(ctx context.Context, urlToSave, alias string) error {
			return repoErr
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Save(context.Background(), "https://example.com", "my-alias")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected db error, got: %v", err)
	}
}

func TestSave_WithoutAlias_GeneratesAlias(t *testing.T) {
	var savedAlias string
	repo := &mockRepo{
		saveFn: func(ctx context.Context, urlToSave, alias string) error {
			savedAlias = alias
			return nil
		},
	}
	gen := &mockGenerator{aliases: []string{"generated1"}}
	svc := NewUrlService(repo, gen, newLogger(), "http://localhost")

	err := svc.Save(context.Background(), "https://example.com", "")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if savedAlias != "http://localhost/generated1" {
		t.Errorf("unexpected alias: %s", savedAlias)
	}
}

func TestSave_WithoutAlias_RetryOnAliasTaken(t *testing.T) {
	callCount := 0
	repo := &mockRepo{
		saveFn: func(ctx context.Context, urlToSave, alias string) error {
			callCount++
			if callCount < 3 {
				return url.ErrAliasTaken
			}
			return nil
		},
	}
	gen := &mockGenerator{aliases: []string{"alias1", "alias2", "alias3"}}
	svc := NewUrlService(repo, gen, newLogger(), "http://localhost")

	err := svc.Save(context.Background(), "https://example.com", "")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 repo calls, got: %d", callCount)
	}
}

func TestSave_WithoutAlias_AllAliasesTaken(t *testing.T) {
	repo := &mockRepo{
		saveFn: func(ctx context.Context, urlToSave, alias string) error {
			return url.ErrAliasTaken
		},
	}
	gen := &mockGenerator{aliases: []string{"a1", "a2", "a3", "a4", "a5"}}
	svc := NewUrlService(repo, gen, newLogger(), "http://localhost")

	// Все 5 попыток провалились — сервис возвращает nil (цикл завершается без ошибки)
	err := svc.Save(context.Background(), "https://example.com", "")
	if err != nil {
		t.Errorf("expected nil, got: %v", err)
	}
}

func TestSave_WithoutAlias_NonAliasError(t *testing.T) {
	repoErr := errors.New("unexpected db error")
	callCount := 0
	repo := &mockRepo{
		saveFn: func(ctx context.Context, urlToSave, alias string) error {
			callCount++
			return repoErr
		},
	}
	gen := &mockGenerator{aliases: []string{"alias1", "alias2"}}
	svc := NewUrlService(repo, gen, newLogger(), "http://localhost")

	err := svc.Save(context.Background(), "https://example.com", "")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 repo call on non-alias error, got: %d", callCount)
	}
}

// --- List tests ---

func TestList_Success(t *testing.T) {
	expected := []url.Url{
		{Id: 1, OriginalUrl: "https://example.com", Alias: "http://localhost/abc"},
		{Id: 2, OriginalUrl: "https://google.com", Alias: "http://localhost/xyz"},
	}
	repo := &mockRepo{
		listFn: func(ctx context.Context) ([]url.Url, error) {
			return expected, nil
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result) != len(expected) {
		t.Errorf("expected %d urls, got %d", len(expected), len(result))
	}
	for i, u := range result {
		if u.Id != expected[i].Id || u.OriginalUrl != expected[i].OriginalUrl {
			t.Errorf("url mismatch at index %d", i)
		}
	}
}

func TestList_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mockRepo{
		listFn: func(ctx context.Context) ([]url.Url, error) {
			return nil, repoErr
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	_, err := svc.List(context.Background())
	if !errors.Is(err, repoErr) {
		t.Errorf("expected db error, got: %v", err)
	}
}

// --- Get tests ---

func TestGet_Success(t *testing.T) {
	expected := url.Url{Id: 42, OriginalUrl: "https://example.com", Alias: "http://localhost/abc"}
	repo := &mockRepo{
		getFn: func(ctx context.Context, id int) (url.Url, error) {
			if id != 42 {
				t.Errorf("expected id 42, got %d", id)
			}
			return expected, nil
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	result, err := svc.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Id != expected.Id || result.OriginalUrl != expected.OriginalUrl {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := &mockRepo{
		getFn: func(ctx context.Context, id int) (url.Url, error) {
			return url.Url{}, url.ErrNotFound
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, url.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

// --- Update tests ---

func TestUpdate_WithAlias_Success(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(ctx context.Context, id int, newUrl, alias string) error {
			if id != 1 {
				t.Errorf("expected id 1, got %d", id)
			}
			if newUrl != "https://new.com" {
				t.Errorf("unexpected url: %s", newUrl)
			}
			if alias != "http://localhost/new-alias" {
				t.Errorf("unexpected alias: %s", alias)
			}
			return nil
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Update(context.Background(), 1, "https://new.com", "new-alias")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestUpdate_WithoutAlias_Success(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(ctx context.Context, id int, newUrl, alias string) error {
			if alias != "" {
				t.Errorf("expected empty alias, got: %s", alias)
			}
			return nil
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Update(context.Background(), 1, "https://new.com", "")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestUpdate_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mockRepo{
		updateFn: func(ctx context.Context, id int, newUrl, alias string) error {
			return repoErr
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Update(context.Background(), 1, "https://new.com", "alias")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected db error, got: %v", err)
	}
}

// --- Delete tests ---

func TestDelete_Success(t *testing.T) {
	repo := &mockRepo{
		deleteFn: func(ctx context.Context, id int) error {
			if id != 5 {
				t.Errorf("expected id 5, got %d", id)
			}
			return nil
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Delete(context.Background(), 5)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestDelete_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	repo := &mockRepo{
		deleteFn: func(ctx context.Context, id int) error {
			return repoErr
		},
	}
	svc := NewUrlService(repo, &mockGenerator{}, newLogger(), "http://localhost")

	err := svc.Delete(context.Background(), 5)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected db error, got: %v", err)
	}
}

// --- BuildShortUrl tests ---

func TestBuildShortUrl(t *testing.T) {
	svc := &urlService{baseUrl: "http://localhost"}

	tests := []struct {
		baseUrl string
		code    string
		want    string
	}{
		{"http://localhost", "abc123", "http://localhost/abc123"},
		{"http://localhost/", "abc123", "http://localhost/abc123"},
		{"http://localhost///", "abc123", "http://localhost/abc123"},
		{"https://short.ly", "xyz", "https://short.ly/xyz"},
	}

	for _, tc := range tests {
		got := svc.BuildShortUrl(tc.baseUrl, tc.code)
		if got != tc.want {
			t.Errorf("BuildShortUrl(%q, %q) = %q, want %q", tc.baseUrl, tc.code, got, tc.want)
		}
	}
}
