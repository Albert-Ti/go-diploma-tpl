package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/repository"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

var (
	ErrLoginAlreadyExists         = errors.New("User already exist")
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
	userID, err := s.repository.Register(ctx, login, hash)
	classifier := repository.NewPostgresErrorClassifier()
	classification := classifier.Classify(err)

	if err != nil {
		if classification == repository.NonRetriable {
			return 0, ErrLoginAlreadyExists
		}
	}
	return userID, nil
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

func (s *Service) CheckAccrualOrder(ctx context.Context, task models.TaskOrder) error {
	if err := s.repository.UpdateStatusOrder(ctx, task.OrderID, models.StatusProcessing); err != nil {
		return err
	}

	urlAccrual := url.URL{
		Scheme: "http",
		Host:   config.Envs.AccrualSystemAddr,
		Path:   "api/orders/" + task.OrderID,
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAccrual.String(), nil)
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			time.Sleep(time.Second * 30)
			continue
		}

		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			var m models.AccrualResp
			if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
				return err
			}

			switch m.Status {
			case "PROCESSED":
				err := s.repository.ProcessedOrder(ctx, task.OrderID, models.StatusProcessed, m.Accrual, task.UserID)
				if err != nil {
					return err
				}
				return nil

			case "INVALID":
				return s.repository.UpdateStatusOrder(ctx, task.OrderID, models.StatusInvalid)

			case "PROCESSING":
				time.Sleep(time.Second * 10)
				continue
			}

		case http.StatusNoContent:
			time.Sleep(time.Second * 10)
			continue

		case http.StatusTooManyRequests:
			retryAfter := res.Header.Get("Retry-After")
			waitTime, _ := strconv.Atoi(retryAfter)
			if waitTime == 0 {
				waitTime = 60
			}
			time.Sleep(time.Second * time.Duration(waitTime))
			continue

		case http.StatusInternalServerError:
			time.Sleep(time.Second * 30)
			continue

		default:
			return fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

	}
}

func (s *Service) GetOrders(ctx context.Context, userID int) ([]models.GetOrdersResp, error) {
	return s.repository.GetOrders(ctx, userID)
}

func (s *Service) GetBalance(ctx context.Context, userID int) (models.GetBalanceResp, error) {
	return s.repository.GetBalance(ctx, userID)
}
