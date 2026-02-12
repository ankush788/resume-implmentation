package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// RegisterRequest represents a user registration request message
type RegisterRequest struct {
	TaskID   string `json:"task_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	RegistrationTopic = "user-registration-requests"
)



// KafkaProducer handles producing messages to Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}


// KafkaConsumer handles consuming messages from Kafka
type KafkaConsumer struct {
	reader *kafka.Reader
}


// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    RegistrationTopic,
		Balancer: &kafka.LeastBytes{},  //Which partition should this message go to -> decision take by balancer
										//LeastBytes strategy : Send message to the partition that currently has the least data.
	}
	return &KafkaProducer{writer: w}
}


// PublishRegistration publishes(storing) a registration request to Kafka 
func (kp *KafkaProducer) PublishRegistration(ctx context.Context, req RegisterRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(req.TaskID), //Kafka uses Key to decide partition. 1) If same TaskID is sent again. -->It will go to the same partition.	
		Value: []byte(data), // actual data 
	}

	return kp.writer.WriteMessages(ctx, message)
}


// Close closes the producer 
//writer -> Closes network connections , Releases memory , Stops background goroutines
func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}


// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    RegistrationTopic,
		GroupID:  "user-registration-group",
		StartOffset: kafka.LastOffset, //Start reading from the latest messages (ignore old ones).
	})
	return &KafkaConsumer{reader: r}
}


// ReadMessage reads a single message from Kafka
func (kc *KafkaConsumer) ReadMessage(ctx context.Context) (*RegisterRequest, error) {
	msg, err := kc.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	var req RegisterRequest
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return nil, err
	}

	return &req, nil
}

// Close closes the consumer
//writer -> Closes network connections , Releases memory , Stops background goroutines
func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}
