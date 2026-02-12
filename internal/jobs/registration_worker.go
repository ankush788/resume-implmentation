package jobs

import (
	"context"
	"log"
	"time"

	"go-user-service/internal/cache"
	"go-user-service/internal/messaging"
	"go-user-service/internal/service"
)

// RegistrationWorker processes registration requests from Kafka
type RegistrationWorker struct {
	consumer *messaging.KafkaConsumer
	service  *service.UserService
	cache    *cache.RedisCache
	stopCh   chan struct{}
}

// NewRegistrationWorker creates a new worker
func NewRegistrationWorker(
	consumer *messaging.KafkaConsumer,
	service *service.UserService,
	cache *cache.RedisCache,
) *RegistrationWorker {
	return &RegistrationWorker{
		consumer: consumer,
		service:  service,
		cache:    cache,
		stopCh:   make(chan struct{}),
	}
}

// Start processing messages from Kafka ( background job)
func (w *RegistrationWorker) Start(ctx context.Context) {
	go func() {
		log.Println("Starting registration worker...")
		for {
			select {   //select is used only with channels. Whichever channel is ready first, run that case.
			case <-w.stopCh:
				log.Println("Stopping registration worker...")
				return
			default:
				// Read message from Kafka with timeout
				readCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				req, err := w.consumer.ReadMessage(readCtx)  // reading message from kafka 
				cancel()

				if err != nil {
					log.Printf("Error reading message from Kafka: %v", err)
					continue
				}

				// Process the registration request
				w.processRegistration(ctx, req)
			}
		}
	}()
}

// processRegistration handles the actual registration and updates status
func (w *RegistrationWorker) processRegistration(ctx context.Context, req *messaging.RegisterRequest) {
	log.Printf("Processing registration for user: %s (taskID: %s)", req.Username, req.TaskID)

	// Update status to "processing"
	if err := w.cache.SetTaskStatus(ctx, req.TaskID, "processing", 1*time.Hour); err != nil {
		log.Printf("Failed to update status to processing: %v", err)
	}

	// Perform registration
	err := w.service.Register(ctx, req.Username, req.Password)

	// Update final status based on result
	if err != nil {
		failureMsg := "failed: " + err.Error()
		log.Printf("Registration failed for taskID %s: %v", req.TaskID, err)
		if err := w.cache.SetTaskStatus(ctx, req.TaskID, failureMsg, 1*time.Hour); err != nil {
			log.Printf("Failed to update status: %v", err)
		}
	} else {
		log.Printf("Registration successful for taskID %s", req.TaskID)
		if err := w.cache.SetTaskStatus(ctx, req.TaskID, "completed", 1*time.Hour); err != nil {
			log.Printf("Failed to update status to completed: %v", err)
		}
	}
}

// Stop gracefully stops the worker
func (w *RegistrationWorker) Stop() {
	close(w.stopCh) //from here signal generate goes to start and background job stop
}
