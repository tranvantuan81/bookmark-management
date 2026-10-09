package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/tranvantuan81/bookmark-management/internal/app/repository"
)

const DefaultCodeLength = 7

var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// CodeGen is the interface for generating a random code
//
//go:generate mockery --name=CodeGen --filename=code_gen.go
type CodeGen interface {
	GenerateCode(length int) (string, error)
}

// ShortenUrl is the interface for the shorten URL service
//
//go:generate mockery --name=ShortenUrl --filename=shortenurl.go
type ShortenUrl interface {
	CreateShortenLink(ctx context.Context, url string, expTime int64) (string, error)
	GetURLFromCode(ctx context.Context, code string) (string, error)
}

type shortenUrl struct {
	r       repository.URLStorage
	codeGen CodeGen
}

// NewShortenUrl creates a new shorten URL service
func NewShortenUrl(r repository.URLStorage, codeGen CodeGen) ShortenUrl {
	return &shortenUrl{r: r, codeGen: codeGen}
}

// CreateShortenLink creates a new shorten link for the given URL

func (s *shortenUrl) CreateShortenLink(ctx context.Context, url string, expTime int64) (string, error) {
	const maxAttempts = 5

	for range maxAttempts {
		code, err := s.codeGen.GenerateCode(DefaultCodeLength)
		if err != nil {
			return "", fmt.Errorf("generate short code: %w", err)
		}

		value, err := s.r.GetURL(ctx, code)
		switch {
		case errors.Is(err, repository.ErrNotFound):

		case err != nil:
			return "", fmt.Errorf("check short code: %w", err)

		case value != "":
			continue
		}

		err = s.r.StoreURL(ctx, code, url, expTime)
		if err != nil {
			return "", fmt.Errorf("store short URL: %w", err)
		}

		return code, nil
	}

	return "", errors.New("failed to generate unique short code")
}

// GetURLFromCode returns the original URL for the given code
func (s *shortenUrl) GetURLFromCode(ctx context.Context, code string) (string, error) {
	return s.r.GetURL(ctx, code)
}

// codeGenImpl is the default implementation of CodeGen
type codeGenImpl struct{}

// NewCodeGen creates a new CodeGen
func NewCodeGen() CodeGen {
	return &codeGenImpl{}
}

// GenerateCode generates a random code of the given length
func (c *codeGenImpl) GenerateCode(length int) (string, error) {
	return GenerateCode(length)
}

// GenerateCode generates a random code of the given length
func GenerateCode(length int) (string, error) {
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
