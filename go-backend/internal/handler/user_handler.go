package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/repository"
	"user-api/internal/service"
)

var validate = validator.New()

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	service service.UserService
	logger  *logger.Logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(service service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(&req); err != nil {
		validationErrors := formatValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	user, err := h.service.CreateUser(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDOB) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date of birth format. Use YYYY-MM-DD",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.service.GetUser(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(&req); err != nil {
		validationErrors := formatValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	user, err := h.service.UpdateUser(c.Context(), int64(id), &req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		if errors.Is(err, service.ErrInvalidDOB) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date of birth format. Use YYYY-MM-DD",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(user)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = h.service.DeleteUser(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Parse pagination query params with defaults
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	// Validate pagination params
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	result, err := h.service.ListUsers(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users",
		})
	}

	return c.JSON(result)
}

// formatValidationErrors converts validator errors to readable format
func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errors[field] = field + " is required"
		case "min":
			errors[field] = field + " must be at least " + err.Param() + " characters"
		case "max":
			errors[field] = field + " must be at most " + err.Param() + " characters"
		case "datetime":
			errors[field] = field + " must be in format " + err.Param()
		default:
			errors[field] = field + " is invalid"
		}
	}

	return errors
}
