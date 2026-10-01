package service

import (
	"context"
	"errors"

	"github.com/tranvantuan81/bookmark-management/internal/repository"
)

const DefaultCodeLength = 7

// ShortenUrl is the interface for the shorten URL service
//
//go:generate mockery --name=ShortenUrl --filename=shortenurl.go
type ShortenUrl interface {
	CreateShortenLink(ctx context.Context, url string, expTime int) (string, error)
}

type shortenUrl struct {
	r       repository.URLStorage
	codeGen GenPass
}

// NewShortenUrl creates a new shorten URL service
func NewShortenUrl(r repository.URLStorage, codeGen GenPass) ShortenUrl {
	return &shortenUrl{r: r, codeGen: codeGen}
}

// CreateShortenLink creates a new shorten link for the given URL
func (s *shortenUrl) CreateShortenLink(ctx context.Context, url string, expTime int) (string, error) {
	for {
		// Generate code
		code, err := s.codeGen.GeneratePassword(DefaultCodeLength)
		if err != nil {
			return "", err
		}

		// Check code exists in storage
		value, err := s.r.GetURL(ctx, code)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return "", err
		}

		if value != "" {
			continue
		}

		err = s.r.StoreURL(ctx, code, url, expTime)
		if err != nil {
			return "", err
		}

		return code, nil
	}
}
