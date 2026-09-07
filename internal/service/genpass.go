package service

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenPass is the interface for the genpass service
//
//go:generate mockery --name=GenPass --filename=genpass.go
type GenPass interface {
	GeneratePassword(length int) (string, error)
}

type genPassService struct {
}

// NewGenPass creates a new genpass service
func NewGenPass() GenPass {
	return &genPassService{}
}

// GeneratePassword generates a new password
func (s *genPassService) GeneratePassword(length int) (string, error) {
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}
