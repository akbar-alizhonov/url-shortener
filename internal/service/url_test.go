package service

import (
	"awesomeProject/internal/repositiries"
	"log/slog"
	"reflect"
	"testing"
)

func TestNewUrlService(t *testing.T) {
	type args struct {
		repo      repositiries.UrlRepository
		generator AliasGenerator
		logger    *slog.Logger
		baseUrl   string
	}
	tests := []struct {
		name string
		args args
		want UrlService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUrlService(tt.args.repo, tt.args.generator, tt.args.logger, tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUrlService() = %v, want %v", got, tt.want)
			}
		})
	}
}
