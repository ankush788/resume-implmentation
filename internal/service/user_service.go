package service

import (
	"fmt"
	"sync"
	"time"

	"go-user-service/internal/errors"
	"go-user-service/internal/models"
	"go-user-service/internal/repository"
	"go-user-service/pkg/utils"
)

type UserService struct {
	repo  *repository.UserRepository
	tasks map[string]string
	mu    sync.RWMutex
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo, make(map[string]string), sync.RWMutex{}}
}

// Register creates a new user with password hashing
func (s *UserService) Register(username, password string) error {
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

	err = s.repo.Create(user)
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
func (s *UserService) Login(username, password string) error {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return errors.ErrInvalidCredentials
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return errors.ErrInvalidCredentials
	}

	return nil
}

// StartRegisterAsync creates a background task that registers a user in a goroutine.
// It returns a task ID that can be used to query status.
func (s *UserService) StartRegisterAsync(username, password string) (string, error) {
	// create a simple task id
	taskID := fmt.Sprintf("%d", time.Now().UnixNano())

	s.mu.Lock()
	s.tasks[taskID] = "pending"
	s.mu.Unlock()

	go func(id, u, p string) {
		// perform registration
		if err := s.Register(u, p); err != nil {
			s.mu.Lock()
			s.tasks[id] = "failed: " + err.Error()
			s.mu.Unlock()
			return
		}
		s.mu.Lock()
		s.tasks[id] = "completed"
		s.mu.Unlock()
	}(taskID, username, password)

	return taskID, nil
}

// GetTaskStatus returns the current status of a background task.
func (s *UserService) GetTaskStatus(taskID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if st, ok := s.tasks[taskID]; ok {
		return st
	}
	return "not_found"
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
