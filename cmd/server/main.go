package main

import (
	"go-user-service/internal/cache"
	"go-user-service/internal/config"
	"go-user-service/internal/db"
	"go-user-service/internal/handler"
	"go-user-service/internal/middleware"
	"go-user-service/internal/models"
	"go-user-service/internal/repository"
	"go-user-service/internal/routes"
	"go-user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.DBUrl)

	// Auto migrate
	//AutoMigrate keeps your database schema in sync with your Go models. if some field in model not in db table, it will add that field to the table. but 
	// it will not delete or modify existing fields to prevent data loss (it throw error).
	database.AutoMigrate(&models.User{})

	// Initialize Redis cache (best-effort)
	rc := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	repo := repository.NewUserRepo(database, rc) // database layer  (database)
	service := service.NewUserService(repo)      // business logic layer (database -> service)
	handler := handler.NewUserHandler(service)   // http layer ( database ->  service -> handler)

    //Gin Engine manages routing and middleware, and for each request it creates a Context  which handles request data, binding, and responses.”
	r := gin.Default()  // r is gin engine pointer
	
	// Add centralized error handling middleware
	r.Use(middleware.ErrorHandlerMiddleware())

	routes.RegisterRoutes(r, handler)

	r.Run(":8080") 
}
