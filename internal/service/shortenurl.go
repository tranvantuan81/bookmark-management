package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/tranvantuan81/bookmark-management/internal/repository"
)

const DefaultCodeLength = 7

var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ShortenUrl is the interface for the shorten URL service
//
//go:generate mockery --name=ShortenUrl --filename=shortenurl.go
type ShortenUrl interface {
	CreateShortenLink(ctx context.Context, url string, expTime int64) (string, error)
	GetURLFromCode(ctx context.Context, code string) (string, error)
}

type shortenUrl struct {
	r repository.URLStorage
}

// NewShortenUrl creates a new shorten URL service
func NewShortenUrl(r repository.URLStorage) ShortenUrl {
	return &shortenUrl{r: r}
}

// CreateShortenLink creates a new shorten link for the given URL
func (s *shortenUrl) CreateShortenLink(ctx context.Context, url string, expTime int64) (string, error) {
	// Generate code
	code, err := GeneratePassword(DefaultCodeLength)
	if err != nil {
		return "", err
	}

	// Check code exists in storage
	value, err := s.r.GetURL(ctx, code)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return "", err
	}

	if value != "" {
		return s.CreateShortenLink(ctx, url, expTime)
	}

	err = s.r.StoreURL(ctx, code, url, expTime)
	if err != nil {
		return "", err
	}

	return code, nil
}

// GetURLFromCode returns the original URL for the given code
func (s *shortenUrl) GetURLFromCode(ctx context.Context, code string) (string, error) {
	return s.r.GetURL(ctx, code)
}

// GeneratePassword generates a random password of the given length
func GeneratePassword(length int) (string, error) {
	password := make([]byte, length)

	for i := range length {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}
