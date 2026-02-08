package url

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrAliasTaken   = errors.New("alias already taken")
	ErrInvalidAlias = errors.New("invalid alias")
)
