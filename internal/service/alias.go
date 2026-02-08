package service

import "math/rand"

type AliasGenerator interface {
	Generate() string
}

type aliasGenerator struct{}

func NewAliasGenerator() AliasGenerator {
	return &aliasGenerator{}
}

func (g *aliasGenerator) Generate() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
