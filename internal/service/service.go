package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/Albert-Ti/go-diploma-tpl/internal/repository"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

var ErrUnauthorized = errors.New("Invalid password")

type Service struct {
	repository repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Register(ctx context.Context, login string, pass string) error {
	salt, err := utils.RandomHash(8)
	if err != nil {
		return err
	}

	hash := hashPassword(salt, pass)
	return s.repository.AddUser(ctx, login, hash)
}

func (s *Service) Login(ctx context.Context, login string, pass string) (int, error) {
	userID, storedHash, err := s.repository.GetUser(ctx, login)

	if err != nil {
		return 0, err
	}
	salt := strings.Split(storedHash, ".")[0]
	hash := hashPassword(salt, pass)

	if storedHash != hash {
		return 0, ErrUnauthorized
	}

	return userID, nil
}

func hashPassword(salt string, pass string) string {
	sum := sha256.Sum256([]byte(pass + salt))
	encStr := base64.StdEncoding.EncodeToString(sum[:])
	return fmt.Sprint(salt, ".", encStr)
}
