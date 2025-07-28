// internal/interfaces/publisher/event_publisher.go
package publisher

type EventPublisher interface {
	Publish(eventType string, payload interface{}) error
}
