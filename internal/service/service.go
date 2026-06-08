package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Albert-Ti/go-diploma-tpl/internal/repository"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

var (
	ErrUnauthorized               = errors.New("Invalid password")
	ErrOrderAlreadyExistsForUser  = errors.New("Order already exists for this user")
	ErrOrderAlreadyExistsForOther = errors.New("Order already exists for other user")
)

type Service struct {
	repository repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) Register(ctx context.Context, login string, pass string) (int, error) {
	salt, err := utils.RandomHash(8)
	if err != nil {
		return 0, err
	}

	hash := utils.HashPassword(salt, pass)
	return s.repository.Registration(ctx, login, hash)
}

func (s *Service) Login(ctx context.Context, login string, pass string) (int, error) {
	userID, storedHash, err := s.repository.GetUser(ctx, login)

	if err != nil {
		return 0, err
	}
	salt := strings.Split(storedHash, ".")[0]
	hash := utils.HashPassword(salt, pass)

	if storedHash != hash {
		return 0, ErrUnauthorized
	}

	return userID, nil
}

func (s *Service) AddOrder(ctx context.Context, order string, userID int) error {
	err := s.repository.CreateOrder(ctx, order, userID)

	classifier := repository.NewPostgresErrorClassifier()
	classification := classifier.Classify(err)

	if err != nil {
		if classification == repository.NonRetriable {
			existingUserID, getErr := s.repository.GetOrderUserID(ctx, order)
			if getErr != nil {
				return fmt.Errorf("Failed to check existing order: %w", getErr)
			}

			if existingUserID == userID {
				return ErrOrderAlreadyExistsForUser
			}
			return ErrOrderAlreadyExistsForOther
		}
	}
	return nil
}

func (s *Service) GetOrders(ctx context.Context, userID int) error {

	return nil
}
