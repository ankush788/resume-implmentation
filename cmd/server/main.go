package main

import (
	"context"
	"log"
	"strings"
	"go-user-service/internal/cache"
	"go-user-service/internal/config"
	"go-user-service/internal/db"
	"go-user-service/internal/handler"
	"go-user-service/internal/jobs"
	"go-user-service/internal/messaging"
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

	
   // kafkabroker  is server Stores messages , Manages topics , Handles partitions
   // cfg.KafkaBrokerAddress = "localhost:9092,localhost:9093" --(split) -->  []string{ "localhost:9092", "localhost:9093" }
   brokers := strings.Split(cfg.KafkaBrokerAddress, ",") 
	kafkaProducer := messaging.NewKafkaProducer(brokers)
	kafkaConsumer := messaging.NewKafkaConsumer(brokers)

	repo := repository.NewUserRepo(database, rc)
	userService := service.NewUserService(repo)

	// Set Kafka producer and Redis cache dependencies on service
	userService.SetDependencies(kafkaProducer, rc)

	// Start registration worker (processes messages from Kafka)
	worker := jobs.NewRegistrationWorker(kafkaConsumer, userService, rc)
	worker.Start(context.Background())
	log.Println("Registration worker started")

	handler := handler.NewUserHandler(userService)

    //Gin Engine manages routing and middleware, and for each request it creates a Context  which handles request data, binding, and responses.”
	r := gin.Default()  // r is gin engine pointer
	
	// Add centralized error handling middleware
	r.Use(middleware.ErrorHandlerMiddleware())

	routes.RegisterRoutes(r, handler)

	r.Run(":8080") 
}
