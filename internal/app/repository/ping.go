package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Ping interface {
	Ping(ctx context.Context) error
}

type ping struct {
	c *redis.Client
}

func NewPing(c *redis.Client) Ping {
	return &ping{
		c: c,
	}
}

func (s *ping) Ping(ctx context.Context) error {
	return s.c.Ping(ctx).Err()
}
