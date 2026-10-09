package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//const expTime = 3600 * time.Second

// ErrNotFound is returned when the requested key does not exist in the storage.
// It lets upper layers handle a missing key without depending on the
// underlying storage implementation.
var ErrNotFound = errors.New("value not found")

// URLStorage is the interface for the URL storage repository
//
//go:generate mockery --name=URLStorage --filename=urlstorage.go
type URLStorage interface {
	StoreURL(ctx context.Context, key, url string, expTime int64) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewURLStorage creates a new URL storage repository
func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

// StoreURL stores a URL in the storage
func (u *urlStorage) StoreURL(ctx context.Context, key, url string, expTime int64) error {
	return u.c.Set(ctx, key, url, time.Duration(expTime)*time.Second).Err()
}

var ErrKeyNotFound = errors.New("key not found")

// GetURL retrieves a URL from the storage
func (u *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	value, err := u.c.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get URL: %w", err)
	}

	return value, nil
}
