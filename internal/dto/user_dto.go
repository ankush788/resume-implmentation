package dto

// package dto is exactly where you define your request/response structures that move between layers.
//DTO = Data Transfer Object

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TaskResponse represents the async task response
type TaskResponse struct {
	TaskID string `json:"task_id"`
}

// TaskStatusResponse represents the task status response
type TaskStatusResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

// FilterParams represents dynamic filtering query parameters
type FilterParams struct {
	Username string `form:"username" binding:"omitempty"`
}

// UserResponse represents a user in the response (without password)
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// PaginatedUsersResponse represents paginated list of users
type PaginatedUsersResponse struct {
	Data       []UserResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
	Limit      int            `json:"limit"`
	Offset     int            `json:"offset"`
	HasMore    bool           `json:"has_more"`
}
