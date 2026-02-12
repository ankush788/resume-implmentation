package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-user-service/internal/cache"
	"go-user-service/internal/errors"
	"go-user-service/internal/messaging"
	"go-user-service/internal/models"
	"go-user-service/internal/repository"
	"go-user-service/pkg/utils"
)

type UserService struct {
	repo     *repository.UserRepository
	producer *messaging.KafkaProducer
	cache    *cache.RedisCache
	mu       sync.RWMutex
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo:  repo,
		cache: nil,
		mu:    sync.RWMutex{},
	}
}

// SetDependencies sets Kafka producer and Redis cache
func (s *UserService) SetDependencies(producer *messaging.KafkaProducer, cache *cache.RedisCache) {
	s.producer = producer
	s.cache = cache
}

// Register creates a new user with password hashing
func (s *UserService) Register(ctx context.Context, username, password string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return errors.NewInternalError(
			errors.ErrHashingPassword.Message,
			err.Error(),
		)
	}

	user := &models.User{
		Username: username,
		Password: hashed,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		// Check if it's a duplicate username error
		if err.Error() == "duplicate username" {
			return errors.ErrUserAlreadyExists
		}
		return errors.NewInternalError(
			errors.ErrDatabaseError.Message,
			err.Error(),
		)
	}

	return nil
}

// Login validates user credentials
func (s *UserService) Login(ctx context.Context, username, password string) error {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return errors.ErrInvalidCredentials
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return errors.ErrInvalidCredentials
	}

	return nil
}

// StartRegisterAsync publishes a registration request to Kafka queue
// Returns a task ID that can be used to poll the status
func (s *UserService) StartRegisterAsync(ctx context.Context, username, password string) (string, error) {
	if s.producer == nil {
		return "", errors.NewInternalError(
			"Kafka producer not initialized",
			"producer is nil",
		)
	}

	// Generate task ID
	taskID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Create registration request message
	req := messaging.RegisterRequest{
		TaskID:   taskID,
		Username: username,
		Password: password,
	}

	// Publish to Kafka (non-blocking, returns immediately)
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // it wait 5 second then timeout/exist
	defer cancel() // it ensure if job end before 5 second then exist (cancel --> safely exist system realse resources)

	if err := s.producer.PublishRegistration(publishCtx, req); err != nil {
		return "", errors.NewInternalError(
			"Failed to publish registration request",
			err.Error(),
		)
	}

	// Set initial status in Redis
	if s.cache != nil {
		cacheCtx, cacheCancel := context.WithTimeout(ctx, 5*time.Second)
		defer cacheCancel()
		s.cache.SetTaskStatus(cacheCtx, taskID, "pending", 1*time.Hour)
	}

	return taskID, nil
}

// GetTaskStatus returns the current status of a background task (async polling API)
func (s *UserService) GetTaskStatus(ctx context.Context, taskID string) string {
	if s.cache == nil {
		return "cache_unavailable"
	}

	cacheCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := s.cache.GetTaskStatus(cacheCtx, taskID)
	if err != nil {
		return "not_found"
	}

	return status
}

// GetAllUsers retrieves all users with pagination and filtering
func (s *UserService) GetAllUsers(limit, offset int, usernameFilter string) ([]models.User, int64, error) {
	users, totalCount, err := s.repo.FindAllWithPagination(limit, offset, usernameFilter)
	if err != nil {
		return nil, 0, errors.NewInternalError(
			errors.ErrDatabaseError.Message,
			err.Error(),
		)
	}
	return users, totalCount, nil
}
