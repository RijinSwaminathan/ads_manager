package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rijin/ads_manager/internal/campaign"
	"github.com/rijin/ads_manager/internal/types"
)

// MessageType represents different types of messages
type MessageType string

const (
	BidRequestMessage  MessageType = "BID_REQUEST"
	BidResponseMessage MessageType = "BID_RESPONSE"
	MetricsMessage     MessageType = "METRICS_UPDATE"
)

// Message represents a generic message structure
type Message struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// Client handles RabbitMQ messaging operations
type Client struct {
	// In a real implementation, this would hold the RabbitMQ connection and channel
	subscribers map[string][]chan Message
	mu          sync.RWMutex
}

// NewClient creates a new RabbitMQ client instance
func NewClient() *Client {
	return &Client{
		subscribers: make(map[string][]chan Message),
	}
}

// PublishBidRequest publishes a bid request message
func (c *Client) PublishBidRequest(ctx context.Context, req *campaign.BidRequest) error {
	msg := Message{
		Type:      BidRequestMessage,
		Payload:   req,
		Timestamp: time.Now(),
	}
	return c.publish("bid_requests", msg)
}

// PublishBidResponse publishes a bid response message
func (c *Client) PublishBidResponse(ctx context.Context, resp *types.BidResponse) error {
	msg := Message{
		Type:      BidResponseMessage,
		Payload:   resp,
		Timestamp: time.Now(),
	}
	return c.publish("bid_responses", msg)
}

// PublishMetrics publishes campaign metrics
func (c *Client) PublishMetrics(ctx context.Context, metrics *campaign.Performance) error {
	msg := Message{
		Type:      MetricsMessage,
		Payload:   metrics,
		Timestamp: time.Now(),
	}
	return c.publish("metrics", msg)
}

// Subscribe subscribes to messages on a specific topic
func (c *Client) Subscribe(topic string) (<-chan Message, error) {
	ch := make(chan Message, 100) // Buffer size of 100 messages

	c.mu.Lock()
	if _, exists := c.subscribers[topic]; !exists {
		c.subscribers[topic] = make([]chan Message, 0)
	}
	c.subscribers[topic] = append(c.subscribers[topic], ch)
	c.mu.Unlock()

	return ch, nil
}

// Unsubscribe removes a subscription
func (c *Client) Unsubscribe(topic string, ch <-chan Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if subs, exists := c.subscribers[topic]; exists {
		for i, sub := range subs {
			if sub == ch {
				c.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(sub)
				break
			}
		}
	}
	return nil
}

// publish sends a message to all subscribers of a topic
func (c *Client) publish(topic string, msg Message) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if subs, exists := c.subscribers[topic]; exists {
		// Marshal message to validate it can be serialized
		_, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}

		// Send to all subscribers
		for _, ch := range subs {
			select {
			case ch <- msg:
				// Message sent successfully
			default:
				// Channel is full, log warning in real implementation
			}
		}
	}
	return nil
}

// Close closes all subscriber channels
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for topic, subs := range c.subscribers {
		for _, ch := range subs {
			close(ch)
		}
		delete(c.subscribers, topic)
	}
	return nil
}
