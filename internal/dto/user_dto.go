package dto

// package dto is exactly where you define your request/response structures that move between layers.
//DTO = Data Transfer Object
//ex :  RegisterRequest represents the registration request payload


// only required is required other are optional (during binding )
// struct me kya kab use hota hai 
// json --> request body
// form --> Query params , URL form data (x-www-form-urlencoded)
// Path param (/users/:id)  --> uri

//default : Sets a value if the field is NOT provided.
//omitempty : Skips validation if field is empty.

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`  
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type TaskResponse struct {
	TaskID string `json:"task_id"`
}

type TaskStatusResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}


type PaginationParams struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}


type FilterParams struct {
	Username string `form:"username" binding:"omitempty"`
}


type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	CreatedAt string `json:"created_at"`
}


type PaginatedUsersResponse struct {
	Data       []UserResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
	Limit      int            `json:"limit"`
	Offset     int            `json:"offset"`
	HasMore    bool           `json:"has_more"`
}
