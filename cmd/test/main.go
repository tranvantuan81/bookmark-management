package main

import (
	"context"

	"github.com/tranvantuan81/bookmark-management/internal/repository"
	"github.com/tranvantuan81/bookmark-management/pkg/redis"
)

func main() {
	rClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	urlStorage := repository.NewURLStorage(rClient)

	err = urlStorage.StoreURL(ctx, "abcdef", "https://google.com", 0)
	if err != nil {
		return
	}

	url, err := urlStorage.GetURL(ctx, "abcdef")
	if err != nil {
		panic(err)
	}

	println(url)
}
