package main

import (
	"fmt"
	"sync"
	"time"
)

// Message represents the fundamental data we pass around.
type Message struct {
	Topic   string
	Payload string
}

// Broker manages subscriptions and publishing of messages.
type Broker struct {
	subscribers map[string][]chan Message
	mutex       sync.RWMutex
}

// NewBroker initializes a new Message Broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan Message),
	}
}

// Subscribe allows a client to listen for messages on a specific topic.
// It returns a channel that will receive the topic's messages.
func (b *Broker) Subscribe(topic string) <-chan Message {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// A buffered channel allows publishers to not block immediately if the receiver is busy
	ch := make(chan Message, 5)
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	
	fmt.Printf("[Broker] New subscriber added to topic: %s\n", topic)
	return ch
}

// Publish broadcasts a message to all subscribers of the topic.
func (b *Broker) Publish(topic string, payload string) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	subs, found := b.subscribers[topic]
	if !found {
		return
	}

	msg := Message{Topic: topic, Payload: payload}
	
	// Send message to all subscribers
	for _, ch := range subs {
		// Non-blocking send, if subscriber channel is full we drop or skip
		// For a reliable broker we might wait or increase buffer size.
		select {
		case ch <- msg:
		default:
			fmt.Printf("[Broker] Warning: Dropped message to a busy subscriber on topic %s\n", topic)
		}
	}
}

func main() {
	fmt.Println("=== Day 49: Go Message Broker (Pub/Sub) ===")

	broker := NewBroker()

	// WaitGroup to wait for our concurrent demo to finish
	var wg sync.WaitGroup

	// Let's create two subscribers for the "news" topic
	newsSub1 := broker.Subscribe("news")
	newsSub2 := broker.Subscribe("news")
	
	// And one subscriber for the "alerts" topic
	alertSub := broker.Subscribe("alerts")

	// Helper function for a subscriber to process messages
	consume := func(id string, ch <-chan Message) {
		defer wg.Done()
		for i := 0; i < 3; i++ { // Each consumer reads 3 messages before exiting
			msg := <-ch
			fmt.Printf("[%s] Received on '%s': %s\n", id, msg.Topic, msg.Payload)
		}
	}

	wg.Add(3)
	go consume("NewsReader-1", newsSub1)
	go consume("NewsReader-2", newsSub2)
	go consume("AlertMonitor", alertSub)

	// Wait a moment for consumers to be ready
	time.Sleep(time.Millisecond * 100)

	// Simulate Publishers
	go func() {
		broker.Publish("news", "Go 1.25 Release Date Announced!")
		time.Sleep(time.Millisecond * 200)
		broker.Publish("news", "New Concurrency Patterns Discovered")
		time.Sleep(time.Millisecond * 200)
		broker.Publish("news", "Top 10 Go Packages of the Year")
	}()

	go func() {
		time.Sleep(time.Millisecond * 50)
		broker.Publish("alerts", "CRITICAL: Server CPU at 99%")
		time.Sleep(time.Millisecond * 300)
		broker.Publish("alerts", "WARN: Memory usage high")
		time.Sleep(time.Millisecond * 200)
		broker.Publish("alerts", "RESOLVED: CPU normalized")
	}()

	// Wait for all consumers to receive their 3 messages
	wg.Wait()
	fmt.Println("=== Simulation Complete ===")
}
