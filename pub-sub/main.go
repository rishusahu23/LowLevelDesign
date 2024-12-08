package main

import (
	"fmt"
	"sync"
)

type Message struct {
	Topic   string
	Payload string
}

type Subscriber interface {
	Consume(msg Message)
	GetId() string
}

type Publisher interface {
	Publish(msg Message)
	AddSubscriber(topic string, sub Subscriber)
	RemoveSubscriber(topic string, subId string)
}

type ConcreteSubscriber struct {
	Id string
}

func (cs *ConcreteSubscriber) Consume(msg Message) {
	fmt.Printf("Subscriber %s received message: %s on topic: %s\n", cs.Id, msg.Payload, msg.Topic)
}

func (cs *ConcreteSubscriber) GetId() string {
	return cs.Id
}

type PubSubService struct {
	subs map[string]map[string]Subscriber
	lock sync.Mutex
}

func NewPubSubService() *PubSubService {
	return &PubSubService{
		subs: make(map[string]map[string]Subscriber),
	}
}

func (ps *PubSubService) Publish(msg Message) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	subs, exist := ps.subs[msg.Topic]
	if !exist {
		fmt.Printf("No subscribers for topic: %s\n", msg.Topic)
		return
	}
	for _, sub := range subs {
		sub.Consume(msg)
	}
}

func (ps *PubSubService) AddSubscriber(topic string, sub Subscriber) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, exists := ps.subs[topic]; !exists {
		ps.subs[topic] = make(map[string]Subscriber)
	}

	ps.subs[topic][sub.GetId()] = sub
	fmt.Printf("Subscriber %s added to topic %s\n", sub.GetId(), topic)
}

func (ps *PubSubService) RemoveSubscriber(topic string, subID string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, exists := ps.subs[topic]; exists {
		delete(ps.subs[topic], subID)
		fmt.Printf("Subscriber %s removed from topic %s\n", subID, topic)
		if len(ps.subs[topic]) == 0 {
			delete(ps.subs, topic)
		}
	}
}

func main() {
	pubSub := NewPubSubService()
	sub1 := &ConcreteSubscriber{
		Id: "id1",
	}
	sub2 := &ConcreteSubscriber{
		Id: "id2",
	}

	pubSub.AddSubscriber("topic1", sub1)
	pubSub.AddSubscriber("topic1", sub2)
	pubSub.AddSubscriber("topic2", sub1)

	pubSub.Publish(Message{Topic: "topic1", Payload: "Hello Topic 1"})
	pubSub.Publish(Message{Topic: "topic2", Payload: "Hello Topic 2"})
	pubSub.Publish(Message{Topic: "topic3", Payload: "No Subscribers Here"})

	pubSub.RemoveSubscriber("topic1", "sub2")
	pubSub.Publish(Message{Topic: "topic1", Payload: "Hello Topic 1 after removal"})
}

//
//import (
//	"fmt"
//	"sync"
//)
//
//// Message represents the data being sent in the pub-sub system.
//type Message struct {
//	Topic   string
//	Payload string
//}
//
//// Subscriber interface for any entity that wants to subscribe to topics.
//type Subscriber interface {
//	Update(msg Message)
//	GetID() string
//}
//
//// Publisher interface for any entity that publishes messages.
//type Publisher interface {
//	Publish(msg Message)
//	AddSubscriber(topic string, sub Subscriber)
//	RemoveSubscriber(topic string, subID string)
//}
//
//// ConcreteSubscriber represents a concrete implementation of a subscriber.
//type ConcreteSubscriber struct {
//	id string
//}
//
//func (cs *ConcreteSubscriber) Update(msg Message) {
//	fmt.Printf("Subscriber %s received message: %s on topic: %s\n", cs.id, msg.Payload, msg.Topic)
//}
//
//func (cs *ConcreteSubscriber) GetID() string {
//	return cs.id
//}
//
//// PubSubService implements the Publisher interface and manages subscriptions.
//type PubSubService struct {
//	subscribers map[string]map[string]Subscriber // topic -> subscriberID -> Subscriber
//	lock        sync.RWMutex
//}
//
//func NewPubSubService() *PubSubService {
//	return &PubSubService{
//		subscribers: make(map[string]map[string]Subscriber),
//	}
//}
//
//func (ps *PubSubService) Publish(msg Message) {
//	ps.lock.RLock()
//	defer ps.lock.RUnlock()
//
//	subs, exists := ps.subscribers[msg.Topic]
//	if !exists {
//		fmt.Printf("No subscribers for topic: %s\n", msg.Topic)
//		return
//	}
//
//	for _, sub := range subs {
//		sub.Update(msg)
//	}
//}
//
//func (ps *PubSubService) AddSubscriber(topic string, sub Subscriber) {
//	ps.lock.Lock()
//	defer ps.lock.Unlock()
//
//	if _, exists := ps.subscribers[topic]; !exists {
//		ps.subscribers[topic] = make(map[string]Subscriber)
//	}
//
//	ps.subscribers[topic][sub.GetID()] = sub
//	fmt.Printf("Subscriber %s added to topic %s\n", sub.GetID(), topic)
//}
//
//func (ps *PubSubService) RemoveSubscriber(topic string, subID string) {
//	ps.lock.Lock()
//	defer ps.lock.Unlock()
//
//	if _, exists := ps.subscribers[topic]; exists {
//		delete(ps.subscribers[topic], subID)
//		fmt.Printf("Subscriber %s removed from topic %s\n", subID, topic)
//		if len(ps.subscribers[topic]) == 0 {
//			delete(ps.subscribers, topic)
//		}
//	}
//}
//
//func main() {
//	pubSub := NewPubSubService()
//
//	sub1 := &ConcreteSubscriber{id: "sub1"}
//	sub2 := &ConcreteSubscriber{id: "sub2"}
//
//	// Add subscribers to topics
//	pubSub.AddSubscriber("topic1", sub1)
//	pubSub.AddSubscriber("topic1", sub2)
//	pubSub.AddSubscriber("topic2", sub1)
//
//	// Publish messages
//	pubSub.Publish(Message{Topic: "topic1", Payload: "Hello Topic 1"})
//	pubSub.Publish(Message{Topic: "topic2", Payload: "Hello Topic 2"})
//	pubSub.Publish(Message{Topic: "topic3", Payload: "No Subscribers Here"})
//
//	// Remove a subscriber
//	pubSub.RemoveSubscriber("topic1", "sub2")
//
//	// Publish again to check the removal effect
//	pubSub.Publish(Message{Topic: "topic1", Payload: "Hello Topic 1 after removal"})
//}
