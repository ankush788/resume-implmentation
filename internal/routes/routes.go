package routes

import (
	"go-user-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handler.UserHandler) {
	api := r.Group("/api") // api group for differentiating api routes
	{
		api.POST("/register", h.Register) 
		api.POST("/login", h.Login)
		api.POST("/register-async", h.RegisterAsync)
		api.GET("/task/:id", h.TaskStatus)
		api.GET("/users", h.GetAllUsers)
	}
}
