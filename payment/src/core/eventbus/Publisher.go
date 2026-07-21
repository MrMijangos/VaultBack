package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ExchangeName es el mismo exchange topic que usan api/ y vault-ai-service
// (ver api/src/core/eventbus/Publisher.go) -- payment/ publica en el mismo
// bus de eventos, no crea uno nuevo.
const ExchangeName = "vault.events"

// Eventos de suscripción -- consumidos, a futuro, por el Notification
// Service (api/) para insertar filas en "notifications" con los subtipos
// suscripcion_activa/renovada/fallida/cancelada/por_vencer. Ese consumidor
// todavía no existe: depende del ALTER a los CHECK de notifications, que
// junto con las tablas de suscripciones/anuncios se deja para el final.
const (
	EventSubscriptionActivated = "subscription.activated"
	EventSubscriptionRenewed   = "subscription.renewed"
	EventSubscriptionFailed    = "subscription.failed"
	EventSubscriptionCanceled  = "subscription.canceled"
	EventSubscriptionExpiring  = "subscription.expiring"
)

type SubscriptionEventPayload struct {
	EventType      string `json:"event_type"`
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
	PlanName       string `json:"plan_name,omitempty"`
}

// EventOrderConfirmed se publica cuando el comprador confirma haber recibido
// el producto y se libera el pago en escrow. api/ (blockchaincertificates)
// todavía no tiene el consumidor para esto -- se agrega cuando se conecte
// el resto de la persistencia, igual que con los eventos subscription.*.
const EventOrderConfirmed = "order.confirmed"

type OrderEventPayload struct {
	EventType string `json:"event_type"`
	OrderID   string `json:"order_id"`
	BuyerID   string `json:"buyer_id"`
	SellerID  string `json:"seller_id"`
	AssetID   string `json:"asset_id"`
}

type Publisher interface {
	PublishSubscriptionEvent(ctx context.Context, payload SubscriptionEventPayload) error
	PublishOrderEvent(ctx context.Context, payload OrderEventPayload) error
}

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

func (p *RabbitMQPublisher) PublishSubscriptionEvent(ctx context.Context, payload SubscriptionEventPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("no se pudo serializar el evento %s: %w", payload.EventType, err)
	}

	err = p.ch.PublishWithContext(ctx, ExchangeName, payload.EventType, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("no se pudo publicar el evento %s: %w", payload.EventType, err)
	}
	return nil
}

func (p *RabbitMQPublisher) PublishOrderEvent(ctx context.Context, payload OrderEventPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("no se pudo serializar el evento %s: %w", payload.EventType, err)
	}

	err = p.ch.PublishWithContext(ctx, ExchangeName, payload.EventType, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("no se pudo publicar el evento %s: %w", payload.EventType, err)
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

// NoopPublisher: mismo criterio que api/ -- si RabbitMQ no está disponible
// al arrancar, el servicio sigue funcionando (crear/cancelar suscripciones
// no se rompe), solo no se notifica al usuario hasta que RabbitMQ vuelva y
// se reinicie el proceso.
type NoopPublisher struct{}

func (NoopPublisher) PublishSubscriptionEvent(_ context.Context, payload SubscriptionEventPayload) error {
	log.Printf("[eventbus] RabbitMQ no disponible, evento %s (subscription_id=%s) no publicado", payload.EventType, payload.SubscriptionID)
	return nil
}

func (NoopPublisher) PublishOrderEvent(_ context.Context, payload OrderEventPayload) error {
	log.Printf("[eventbus] RabbitMQ no disponible, evento %s (order_id=%s) no publicado", payload.EventType, payload.OrderID)
	return nil
}
