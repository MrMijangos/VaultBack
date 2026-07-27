package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ExchangeName es el exchange topic que comparten el backend Go y
// vault-ai-service (Python). Debe coincidir con rabbitmq_exchange en
// src/infrastructure/config/settings.py del servicio de NLP/ML.
const ExchangeName = "vault.events"

// Publisher envía eventos de dominio (post.created, comment.created,
// review.created, asset.updated) para que vault-ai-service corra NLP/ML
// sobre el contenido nuevo y, en el caso de post/comment/review, publique
// nlp.analyzed de vuelta (ver Consumer.go).
type Publisher interface {
	Publish(ctx context.Context, eventType string, userID string, sourceID string, text *string) error

	// PublishEvent publica un payload propio del evento (no el shape fijo de
	// Publish) -- lo usan asset.created y maintenance.registered, que
	// vault-blockchain consume con campos que Publish no contempla (name,
	// category, created_at / maintenance_id, type). El payload debe incluir
	// su propio "event_type" via json tag si se quiere que viaje en el body.
	PublishEvent(ctx context.Context, eventType string, payload any) error
}

type eventPayload struct {
	EventType string  `json:"event_type"`
	UserID    string  `json:"user_id"`
	SourceID  string  `json:"source_id"`
	Text      *string `json:"text,omitempty"`
}

// RabbitMQPublisher es la implementación real, usada cuando hay conexión
// a RabbitMQ disponible.
type RabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("no se pudo abrir el canal de RabbitMQ: %w", err)
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("no se pudo declarar el exchange %s: %w", ExchangeName, err)
	}

	return &RabbitMQPublisher{conn: conn, ch: ch}, nil
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, eventType string, userID string, sourceID string, text *string) error {
	body, err := json.Marshal(eventPayload{
		EventType: eventType,
		UserID:    userID,
		SourceID:  sourceID,
		Text:      text,
	})
	if err != nil {
		return fmt.Errorf("no se pudo serializar el evento %s: %w", eventType, err)
	}

	// routing_key = event_type: vault-ai-service enlaza su cola exactamente
	// a estas routing keys (post.created, comment.created, review.created,
	// asset.updated).
	err = p.ch.PublishWithContext(ctx, ExchangeName, eventType, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("no se pudo publicar el evento %s: %w", eventType, err)
	}
	return nil
}

// PublishEvent serializa el payload tal cual (a diferencia de Publish, no lo
// envuelve en eventPayload) -- quien llama controla el shape completo.
func (p *RabbitMQPublisher) PublishEvent(ctx context.Context, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("no se pudo serializar el evento %s: %w", eventType, err)
	}

	err = p.ch.PublishWithContext(ctx, ExchangeName, eventType, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("no se pudo publicar el evento %s: %w", eventType, err)
	}
	return nil
}

func (p *RabbitMQPublisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

// NoopPublisher se usa cuando RabbitMQ no está disponible al arrancar: la
// API sigue funcionando con normalidad (crear posts/comentarios/reseñas/
// activos no se rompe), simplemente no se dispara NLP/ML para el contenido
// nuevo hasta que RabbitMQ vuelva a estar disponible y se reinicie la API.
type NoopPublisher struct{}

func (NoopPublisher) Publish(_ context.Context, eventType string, _ string, sourceID string, _ *string) error {
	log.Printf("[eventbus] RabbitMQ no disponible, evento %s (source_id=%s) no publicado", eventType, sourceID)
	return nil
}

func (NoopPublisher) PublishEvent(_ context.Context, eventType string, _ any) error {
	log.Printf("[eventbus] RabbitMQ no disponible, evento %s no publicado", eventType)
	return nil
}
