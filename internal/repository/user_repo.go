package repository

import (
	"context"
	"go-user-service/internal/cache"
	"go-user-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
 	db *gorm.DB
	cache *cache.RedisCache
}

// NewUserRepo initializes a repository with optional Redis cache (pass nil to disable caching)
func NewUserRepo(db *gorm.DB, rc *cache.RedisCache) *UserRepository {
	return &UserRepository{db: db, cache: rc}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	// set cache (best-effort)
	if r.cache != nil {
		// ignore cache errors; DB write is source of truth
		_ = r.cache.SetUser(ctx, user, 24*time.Hour)
	}

	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	// Try cache first
	if r.cache != nil {
		if u, err := r.cache.GetUser(ctx, username); err == nil && u != nil {
			return u, nil
		}
		// on error fallthrough to DB
	}

	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	// populate cache (best-effort)
	if r.cache != nil {
		_ = r.cache.SetUser(ctx, &user, 24*time.Hour)
	}

	return &user, nil
}

// FindAllWithPagination retrieves users with pagination and filtering
func (r *UserRepository) FindAllWithPagination(limit, offset int, usernameFilter string) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64

	// Set default pagination values
	if limit <= 0 || limit > 100 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0 // default offset
	}

	// Build query
	query := r.db

	// Apply filter if provided
	if usernameFilter != "" {
		query = query.Where("username LIKE ?", "%"+usernameFilter+"%")
	}

	// Get total count
	if err := query.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch data
	err := query.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
