package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Message represents the data being sent in the pub-sub system.
type Message struct {
	Topic   string
	Payload string
}

// Subscriber interface for any entity that wants to subscribe to topics.
type Subscriber interface {
	Process(msg Message)
	GetID() string
}

// ConcreteSubscriber represents a subscriber in the distributed system.
type ConcreteSubscriber struct {
	id string
}

func (cs *ConcreteSubscriber) Process(msg Message) {
	fmt.Printf("Subscriber %s received message: %s on topic: %s\n", cs.id, msg.Payload, msg.Topic)
}

func (cs *ConcreteSubscriber) GetID() string {
	return cs.id
}

// RedisPubSubService is the distributed implementation of the pub-sub system.
type RedisPubSubService struct {
	redisClient *redis.Client
	lock        sync.RWMutex
}

func NewRedisPubSubService(redisAddr string) *RedisPubSubService {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &RedisPubSubService{redisClient: client}
}

func (rps *RedisPubSubService) Publish(topic string, payload string) {
	err := rps.redisClient.Publish(context.Background(), topic, payload).Err()
	if err != nil {
		log.Printf("Failed to publish message: %v\n", err)
	}
}

func (rps *RedisPubSubService) Subscribe(topic string, sub Subscriber) {
	go func() {
		pubsub := rps.redisClient.Subscribe(context.Background(), topic)
		defer pubsub.Close()

		// Listen for messages
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				log.Printf("Error receiving message: %v\n", err)
				continue
			}
			sub.Process(Message{Topic: msg.Channel, Payload: msg.Payload})
		}
	}()
}

func main() {
	// Initialize Redis Pub-Sub service
	redisAddr := "localhost:6379" // Redis server address
	pubSubService := NewRedisPubSubService(redisAddr)

	// Create subscribers
	sub1 := &ConcreteSubscriber{id: "sub1"}
	sub2 := &ConcreteSubscriber{id: "sub2"}

	// Subscribe to topics
	pubSubService.Subscribe("topic1", sub1)
	pubSubService.Subscribe("topic1", sub2)

	// Simulate publishers
	go func() {
		for i := 1; i <= 5; i++ {
			pubSubService.Publish("topic1", fmt.Sprintf("Message %d to topic1", i))
			time.Sleep(1 * time.Second)
		}
	}()

	// Keep the main function running to listen for messages
	select {}
}
