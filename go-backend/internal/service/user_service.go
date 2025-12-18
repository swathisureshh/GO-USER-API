package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/repository"
)

var (
	ErrInvalidDOB = errors.New("invalid date of birth format")
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUser(ctx context.Context, id int64) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, page, pageSize int) (*models.PaginatedResponse, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *logger.Logger
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.UserRepository, logger *logger.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	s.logger.Info("Creating new user", zap.String("name", req.Name))

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("Invalid DOB format", zap.Error(err))
		return nil, ErrInvalidDOB
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("User created successfully", zap.Int64("user_id", user.ID))

	response := user.ToResponse(false) // Don't include age in create response
	return &response, nil
}

// GetUser retrieves a user by ID with calculated age
func (s *userService) GetUser(ctx context.Context, id int64) (*models.UserResponse, error) {
	s.logger.Debug("Fetching user", zap.Int64("user_id", id))

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Warn("User not found", zap.Int64("user_id", id))
			return nil, err
		}
		s.logger.Error("Failed to fetch user", zap.Error(err))
		return nil, err
	}

	response := user.ToResponse(true) // Include age in get response
	return &response, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, id int64, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	s.logger.Info("Updating user", zap.Int64("user_id", id))

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("Invalid DOB format", zap.Error(err))
		return nil, ErrInvalidDOB
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Warn("User not found for update", zap.Int64("user_id", id))
			return nil, err
		}
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("User updated successfully", zap.Int64("user_id", user.ID))

	response := user.ToResponse(false) // Don't include age in update response
	return &response, nil
}

// DeleteUser removes a user
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	s.logger.Info("Deleting user", zap.Int64("user_id", id))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Warn("User not found for deletion", zap.Int64("user_id", id))
			return err
		}
		s.logger.Error("Failed to delete user", zap.Error(err))
		return err
	}

	s.logger.Info("User deleted successfully", zap.Int64("user_id", id))
	return nil
}

// ListUsers retrieves paginated list of users
func (s *userService) ListUsers(ctx context.Context, page, pageSize int) (*models.PaginatedResponse, error) {
	s.logger.Debug("Listing users", zap.Int("page", page), zap.Int("page_size", pageSize))

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get users
	users, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	// Get total count
	totalCount, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.Error("Failed to count users", zap.Error(err))
		return nil, err
	}

	// Convert to response with age
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse(true)) // Include age in list response
	}

	// Calculate total pages
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return &models.PaginatedResponse{
		Data:       responses,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}
