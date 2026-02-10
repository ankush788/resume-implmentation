package handler

import (
	"go-user-service/internal/dto"
	appErrors "go-user-service/internal/errors"
	"go-user-service/internal/middleware"
	"go-user-service/internal/service"
	"go-user-service/internal/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s}
}

// Register handles user registration with request validation
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate JSON structure
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondWithError(c, appErrors.ErrInvalidRequest)
		return
	}

	// Validate request fields
	if err := validators.ValidateRegisterRequest(req.Username, req.Password); err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	// Register user through service
	if err := h.service.Register(req.Username, req.Password); err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	middleware.RespondWithSuccess(c, http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// Login handles user login with request validation
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind and validate JSON structure
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondWithError(c, appErrors.ErrInvalidRequest)
		return
	}

	// Validate request fields
	if err := validators.ValidateLoginRequest(req.Username, req.Password); err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	// Perform login through service
	if err := h.service.Login(req.Username, req.Password); err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	middleware.RespondWithSuccess(c, http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Login successful",
	})
}

// RegisterAsync handles asynchronous user registration with validation
func (h *UserHandler) RegisterAsync(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate JSON structure
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.RespondWithError(c, appErrors.ErrInvalidRequest)
		return
	}

	// Validate request fields
	if err := validators.ValidateRegisterRequest(req.Username, req.Password); err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	// Start async registration
	taskID, err := h.service.StartRegisterAsync(req.Username, req.Password)
	if err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	middleware.RespondWithSuccess(c, http.StatusAccepted, dto.TaskResponse{
		TaskID: taskID,
	})
}

// TaskStatus retrieves the status of an async task
func (h *UserHandler) TaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		middleware.RespondWithError(c, appErrors.ErrTaskNotFound)
		return
	}

	status := h.service.GetTaskStatus(taskID)
	if status == "not_found" {
		middleware.RespondWithError(c, appErrors.ErrTaskNotFound)
		return
	}

	middleware.RespondWithSuccess(c, http.StatusOK, dto.TaskStatusResponse{
		TaskID: taskID,
		Status: status,
	})
}

// GetAllUsers handles retrieving users with pagination and filtering
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var pagination dto.PaginationParams
	var filter dto.FilterParams

	// Bind pagination parameters
	if err := c.ShouldBindQuery(&pagination); err != nil {
		middleware.RespondWithError(c, appErrors.ErrInvalidRequest)
		return
	}

	// Bind filter parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		middleware.RespondWithError(c, appErrors.ErrInvalidRequest)
		return
	}

	// Set default pagination values
	limit := pagination.Limit
	if limit == 0 {
		limit = 10
	}
	offset := pagination.Offset

	// Fetch users from service
	users, totalCount, err := h.service.GetAllUsers(limit, offset, filter.Username)
	if err != nil {
		middleware.RespondWithError(c, err)
		return
	}

	// Convert models to response DTOs (exclude passwords)
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	// Calculate hasMore
	hasMore := (int64(offset+limit) < totalCount)

	response := dto.PaginatedUsersResponse{
		Data:       userResponses,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
	}

	middleware.RespondWithSuccess(c, http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    response,
	})
}
