package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	nlpAnalyzedQueue      = "vault-community-service.nlp-analyzed"
	nlpAnalyzedRoutingKey = "nlp.analyzed"

	blockchainConfirmedQueue      = "vault-api.blockchain-confirmed"
	blockchainConfirmedRoutingKey = "blockchain.confirmed"

	// "subscription.*" -- un solo binding con wildcard cubre los cinco
	// event_type que payment/ publica (activated/renewed/failed/canceled/
	// expiring, ver payment/src/core/eventbus/Publisher.go) sin declarar
	// una cola por cada uno.
	subscriptionEventsQueue      = "vault-api.subscription-events"
	subscriptionEventsRoutingKey = "subscription.*"

	orderConfirmedQueue      = "vault-api.order-confirmed"
	orderConfirmedRoutingKey = "order.confirmed"

	orderCreatedQueue      = "vault-api.order-created"
	orderCreatedRoutingKey = "order.created"

	orderShippedQueue      = "vault-api.order-shipped"
	orderShippedRoutingKey = "order.shipped"
)

// nlpAnalyzedPayload es el resultado que vault-ai-service publica tras
// analizar un post/comentario/reseña. Debe coincidir con lo que arma
// RabbitMQPublisher.publish_nlp_analyzed en el servicio de Python
// (source_id, source_type + NLPAnalyzeResponseDTO). entities/topics se
// ignoran porque no hay columna donde guardarlos todavía.
type nlpAnalyzedPayload struct {
	SourceID       string  `json:"source_id"`
	SourceType     string  `json:"source_type"`
	SentimentScore float64 `json:"sentiment_score"`
	SentimentLabel string  `json:"sentiment_label"`
	ToxicityScore  float64 `json:"toxicity_score"`
	IsToxic        bool    `json:"is_toxic"`
}

// StartNLPAnalyzedConsumer escucha nlp.analyzed y actualiza
// sentiment_score/sentiment_label/toxicity_score/is_visible en la tabla
// correspondiente (posts, reviews o comments — comments no tiene sentiment).
// Si RabbitMQ no está disponible no bloquea el arranque de la API: solo
// registra el problema y sigue sin consumir.
func StartNLPAnalyzedConsumer(url string, pool *pgxpool.Pool) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("[eventbus] no se pudo conectar a RabbitMQ para consumir nlp.analyzed: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[eventbus] no se pudo abrir canal para consumir nlp.analyzed: %v", err)
		conn.Close()
		return
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo declarar el exchange %s: %v", ExchangeName, err)
		ch.Close()
		conn.Close()
		return
	}

	queue, err := ch.QueueDeclare(nlpAnalyzedQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo declarar la cola %s: %v", nlpAnalyzedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	if err := ch.QueueBind(queue.Name, nlpAnalyzedRoutingKey, ExchangeName, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo enlazar la cola %s: %v", nlpAnalyzedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo iniciar el consumo de %s: %v", nlpAnalyzedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	log.Printf("[eventbus] escuchando %s en la cola %s", nlpAnalyzedRoutingKey, nlpAnalyzedQueue)

	go func() {
		for msg := range msgs {
			handleNLPAnalyzed(pool, msg)
		}
	}()
}

func handleNLPAnalyzed(pool *pgxpool.Pool, msg amqp.Delivery) {
	defer msg.Ack(false)

	var payload nlpAnalyzedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[eventbus] mensaje nlp.analyzed inválido: %v", err)
		return
	}

	isVisible := !payload.IsToxic
	ctx := context.Background()

	var query string
	var args []any

	switch payload.SourceType {
	case "post":
		query = `UPDATE posts SET sentiment_score = $1, sentiment_label = $2, toxicity_score = $3, is_visible = $4 WHERE id = $5`
		args = []any{payload.SentimentScore, payload.SentimentLabel, payload.ToxicityScore, isVisible, payload.SourceID}
	case "review":
		query = `UPDATE reviews SET sentiment_score = $1, sentiment_label = $2, toxicity_score = $3, is_visible = $4 WHERE id = $5`
		args = []any{payload.SentimentScore, payload.SentimentLabel, payload.ToxicityScore, isVisible, payload.SourceID}
	case "comment":
		query = `UPDATE comments SET toxicity_score = $1, is_visible = $2 WHERE id = $3`
		args = []any{payload.ToxicityScore, isVisible, payload.SourceID}
	default:
		log.Printf("[eventbus] source_type desconocido en nlp.analyzed: %q", payload.SourceType)
		return
	}

	if _, err := pool.Exec(ctx, query, args...); err != nil {
		log.Printf("[eventbus] no se pudo actualizar %s %s con el resultado de NLP: %v", payload.SourceType, payload.SourceID, err)
	}
}

// blockchainConfirmedPayload es lo que publica vault-blockchain al finalizar
// una certificación en Vara (certify_routes.js: certifyInBackground). action
// es REGISTERED/MAINTAINED/RESTORED, se ignora acá -- el título/cuerpo de la
// notificación es genérico para los tres casos.
type blockchainConfirmedPayload struct {
	AssetID     string `json:"asset_id"`
	OwnerID     string `json:"owner_id"`
	TxID        string `json:"tx_id"`
	ConfirmedAt string `json:"confirmed_at"`
	Action      string `json:"action"`
}

// StartBlockchainConfirmedConsumer escucha blockchain.confirmed y crea una
// notificación (type=blockchain, subtype=asset_verificado) para el dueño del
// activo -- el trigger de la tabla notifications (init.sql) ya se encarga de
// avisarle a realtime/ por pg_notify para el push por WebSocket. Mismo
// criterio que StartNLPAnalyzedConsumer: si RabbitMQ no está disponible, no
// bloquea el arranque de la API.
func StartBlockchainConfirmedConsumer(url string, pool *pgxpool.Pool) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("[eventbus] no se pudo conectar a RabbitMQ para consumir blockchain.confirmed: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[eventbus] no se pudo abrir canal para consumir blockchain.confirmed: %v", err)
		conn.Close()
		return
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo declarar el exchange %s: %v", ExchangeName, err)
		ch.Close()
		conn.Close()
		return
	}

	queue, err := ch.QueueDeclare(blockchainConfirmedQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo declarar la cola %s: %v", blockchainConfirmedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	if err := ch.QueueBind(queue.Name, blockchainConfirmedRoutingKey, ExchangeName, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo enlazar la cola %s: %v", blockchainConfirmedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo iniciar el consumo de %s: %v", blockchainConfirmedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	log.Printf("[eventbus] escuchando %s en la cola %s", blockchainConfirmedRoutingKey, blockchainConfirmedQueue)

	go func() {
		for msg := range msgs {
			handleBlockchainConfirmed(pool, msg)
		}
	}()
}

func handleBlockchainConfirmed(pool *pgxpool.Pool, msg amqp.Delivery) {
	defer msg.Ack(false)

	var payload blockchainConfirmedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[eventbus] mensaje blockchain.confirmed inválido: %v", err)
		return
	}

	data, err := json.Marshal(map[string]string{"asset_id": payload.AssetID, "tx_id": payload.TxID})
	if err != nil {
		log.Printf("[eventbus] no se pudo serializar data de la notificacion blockchain.confirmed: %v", err)
		return
	}

	const query = `
		INSERT INTO notifications (user_id, type, subtype, title, body, data)
		VALUES ($1, 'blockchain', 'asset_verificado', $2, $3, $4)
	`
	const title = "Activo verificado en Vara"
	const body = "Tu activo fue certificado en la blockchain. Ya puedes ver su certificado."

	if _, err := pool.Exec(context.Background(), query, payload.OwnerID, title, body, data); err != nil {
		log.Printf("[eventbus] no se pudo crear la notificacion para el activo %s: %v", payload.AssetID, err)
	}
}

// subscriptionEventPayload -- debe coincidir con
// eventbus.SubscriptionEventPayload en payment/src/core/eventbus/Publisher.go.
type subscriptionEventPayload struct {
	EventType      string `json:"event_type"`
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
}

// StartSubscriptionEventsConsumer escucha subscription.* (payment/publica
// activated/renewed/failed/canceled/expiring al crear, renovar, fallar o
// cancelar una suscripción -- HandleStripeWebhookUseCase.go y
// CreateSubscriptionUseCase.go) y crea la notificación correspondiente.
// Antes de esto el evento se publicaba pero nadie lo consumía: un usuario
// al que Stripe le cancelaba la suscripción (o le fallaba la renovación)
// no se enteraba por ningún lado dentro de la app.
func StartSubscriptionEventsConsumer(url string, pool *pgxpool.Pool) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("[eventbus] no se pudo conectar a RabbitMQ para consumir subscription.*: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[eventbus] no se pudo abrir canal para consumir subscription.*: %v", err)
		conn.Close()
		return
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo declarar el exchange %s: %v", ExchangeName, err)
		ch.Close()
		conn.Close()
		return
	}

	queue, err := ch.QueueDeclare(subscriptionEventsQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo declarar la cola %s: %v", subscriptionEventsQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	if err := ch.QueueBind(queue.Name, subscriptionEventsRoutingKey, ExchangeName, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo enlazar la cola %s: %v", subscriptionEventsQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo iniciar el consumo de %s: %v", subscriptionEventsQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	log.Printf("[eventbus] escuchando %s en la cola %s", subscriptionEventsRoutingKey, subscriptionEventsQueue)

	go func() {
		for msg := range msgs {
			handleSubscriptionEvent(pool, msg)
		}
	}()
}

// subscriptionNotificationCopy traduce el event_type de payment/ al
// subtype/título/cuerpo de la notificación. event_type desconocido
// devuelve subtype="" -- el llamador lo descarta en vez de violar el
// CHECK de notifications.subtype.
func subscriptionNotificationCopy(eventType string) (subtype, title, body string) {
	switch eventType {
	case "subscription.activated":
		return "suscripcion_activa", "Suscripción activada", "Tu suscripción ya está activa."
	case "subscription.renewed":
		return "suscripcion_renovada", "Suscripción renovada", "Tu suscripción se renovó correctamente."
	case "subscription.failed":
		return "suscripcion_fallida", "No se pudo cobrar tu suscripción",
			"Actualiza tu método de pago para no perder tus anuncios activos."
	case "subscription.canceled":
		return "suscripcion_cancelada", "Suscripción cancelada", "Tu suscripción y tus anuncios activos se cancelaron."
	case "subscription.expiring":
		return "suscripcion_por_vencer", "Tu suscripción está por vencer", "Renueva pronto para no perder tus beneficios."
	default:
		return "", "", ""
	}
}

func handleSubscriptionEvent(pool *pgxpool.Pool, msg amqp.Delivery) {
	defer msg.Ack(false)

	var payload subscriptionEventPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[eventbus] mensaje de evento de suscripción inválido: %v", err)
		return
	}

	subtype, title, body := subscriptionNotificationCopy(payload.EventType)
	if subtype == "" {
		log.Printf("[eventbus] event_type de suscripción desconocido: %q", payload.EventType)
		return
	}

	data, err := json.Marshal(map[string]string{"subscription_id": payload.SubscriptionID})
	if err != nil {
		log.Printf("[eventbus] no se pudo serializar data de la notificacion de suscripcion: %v", err)
		return
	}

	const query = `
		INSERT INTO notifications (user_id, type, subtype, title, body, data)
		VALUES ($1, 'suscripcion', $2, $3, $4, $5)
	`
	if _, err := pool.Exec(context.Background(), query, payload.UserID, subtype, title, body, data); err != nil {
		log.Printf("[eventbus] no se pudo crear la notificacion de suscripcion %s: %v", payload.EventType, err)
	}
}

// orderConfirmedPayload -- debe coincidir con eventbus.OrderEventPayload en
// payment/src/core/eventbus/Publisher.go.
type orderConfirmedPayload struct {
	OrderID  string `json:"order_id"`
	BuyerID  string `json:"buyer_id"`
	SellerID string `json:"seller_id"`
	AssetID  string `json:"asset_id"`
}

// StartOrderConfirmedConsumer escucha order.confirmed (payment/ lo publica
// en ConfirmOrderUseCase.go al liberar el escrow) y transfiere la propiedad
// real del activo al comprador -- antes de esto solo quedaba el certificado
// en Vara (ver vault-blockchain/rabbitmq_consumer.js), pero el activo
// seguía apareciendo en el perfil del vendedor original en la propia app.
func StartOrderConfirmedConsumer(url string, pool *pgxpool.Pool) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("[eventbus] no se pudo conectar a RabbitMQ para consumir order.confirmed: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[eventbus] no se pudo abrir canal para consumir order.confirmed: %v", err)
		conn.Close()
		return
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo declarar el exchange %s: %v", ExchangeName, err)
		ch.Close()
		conn.Close()
		return
	}

	queue, err := ch.QueueDeclare(orderConfirmedQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo declarar la cola %s: %v", orderConfirmedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	if err := ch.QueueBind(queue.Name, orderConfirmedRoutingKey, ExchangeName, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo enlazar la cola %s: %v", orderConfirmedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo iniciar el consumo de %s: %v", orderConfirmedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	log.Printf("[eventbus] escuchando %s en la cola %s", orderConfirmedRoutingKey, orderConfirmedQueue)

	go func() {
		for msg := range msgs {
			handleOrderConfirmed(pool, msg)
		}
	}()
}

func handleOrderConfirmed(pool *pgxpool.Pool, msg amqp.Delivery) {
	defer msg.Ack(false)

	var payload orderConfirmedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[eventbus] mensaje order.confirmed inválido: %v", err)
		return
	}

	// El activo deja de estar "en venta" bajo el nuevo dueño -- se limpia
	// junto con la transferencia en la misma sentencia, no en un paso aparte.
	const transferQuery = `
		UPDATE assets
		SET user_id = $1, is_for_sale = false, sale_price = NULL, sale_description = NULL
		WHERE id = $2
		RETURNING name
	`
	var assetName string
	if err := pool.QueryRow(context.Background(), transferQuery, payload.BuyerID, payload.AssetID).Scan(&assetName); err != nil {
		log.Printf("[eventbus] no se pudo transferir la propiedad del activo %s: %v", payload.AssetID, err)
		return
	}

	data, err := json.Marshal(map[string]string{"asset_id": payload.AssetID, "order_id": payload.OrderID})
	if err != nil {
		log.Printf("[eventbus] no se pudo serializar data de la notificacion de compra: %v", err)
		return
	}

	const notifyQuery = `
		INSERT INTO notifications (user_id, type, subtype, title, body, data)
		VALUES ($1, 'venta', 'nueva_compra', $2, $3, $4)
	`
	title := "Compra confirmada"
	body := fmt.Sprintf("Ya eres el dueño de %s. Puedes ver su certificado en Vara desde su detalle.", assetName)
	if _, err := pool.Exec(context.Background(), notifyQuery, payload.BuyerID, title, body, data); err != nil {
		log.Printf("[eventbus] no se pudo crear la notificacion de compra para el activo %s: %v", payload.AssetID, err)
	}
}

// StartOrderCreatedConsumer escucha order.created (payment/ lo publica en
// CreateOrderUseCase.go justo después de cobrar al comprador) y avisa al
// vendedor de la venta nueva -- antes de esto el vendedor solo se enteraba
// si entraba a revisar manualmente, no había notificación ni forma de ver
// sus pedidos (ver ListMySalesUseCase.go en payment/).
func StartOrderCreatedConsumer(url string, pool *pgxpool.Pool) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("[eventbus] no se pudo conectar a RabbitMQ para consumir order.created: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[eventbus] no se pudo abrir canal para consumir order.created: %v", err)
		conn.Close()
		return
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo declarar el exchange %s: %v", ExchangeName, err)
		ch.Close()
		conn.Close()
		return
	}

	queue, err := ch.QueueDeclare(orderCreatedQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo declarar la cola %s: %v", orderCreatedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	if err := ch.QueueBind(queue.Name, orderCreatedRoutingKey, ExchangeName, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo enlazar la cola %s: %v", orderCreatedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo iniciar el consumo de %s: %v", orderCreatedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	log.Printf("[eventbus] escuchando %s en la cola %s", orderCreatedRoutingKey, orderCreatedQueue)

	go func() {
		for msg := range msgs {
			handleOrderCreated(pool, msg)
		}
	}()
}

func handleOrderCreated(pool *pgxpool.Pool, msg amqp.Delivery) {
	defer msg.Ack(false)

	var payload orderConfirmedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[eventbus] mensaje order.created inválido: %v", err)
		return
	}

	// El activo deja de listarse en el Shop en cuanto se paga (no hasta que
	// el comprador confirme recibido) -- antes seguía apareciendo "en venta"
	// durante todo el tránsito (retenido/enviado), pudiendo venderse dos
	// veces. La transferencia real de dueño sigue pasando en
	// handleOrderConfirmed, esto solo lo saca del catálogo.
	const hideQuery = `UPDATE assets SET is_for_sale = false WHERE id = $1`
	if _, err := pool.Exec(context.Background(), hideQuery, payload.AssetID); err != nil {
		log.Printf("[eventbus] no se pudo ocultar el activo %s del catálogo: %v", payload.AssetID, err)
	}

	data, err := json.Marshal(map[string]string{"asset_id": payload.AssetID, "order_id": payload.OrderID})
	if err != nil {
		log.Printf("[eventbus] no se pudo serializar data de la notificacion de venta nueva: %v", err)
		return
	}

	const query = `
		INSERT INTO notifications (user_id, type, subtype, title, body, data)
		VALUES ($1, 'venta', 'pedido_recibido', $2, $3, $4)
	`
	const title = "Tienes una venta nueva"
	const body = "Un comprador ya pagó tu producto. Márcalo como enviado cuando lo despaches."
	if _, err := pool.Exec(context.Background(), query, payload.SellerID, title, body, data); err != nil {
		log.Printf("[eventbus] no se pudo crear la notificacion de venta nueva para la orden %s: %v", payload.OrderID, err)
	}
}

// StartOrderShippedConsumer escucha order.shipped (payment/ lo publica en
// ShipOrderUseCase.go cuando el vendedor marca el pedido como enviado) y
// avisa al comprador -- antes de esto no existía el paso "enviado" ni una
// forma de que el comprador supiera que ya podía esperar su pedido.
func StartOrderShippedConsumer(url string, pool *pgxpool.Pool) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("[eventbus] no se pudo conectar a RabbitMQ para consumir order.shipped: %v", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[eventbus] no se pudo abrir canal para consumir order.shipped: %v", err)
		conn.Close()
		return
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo declarar el exchange %s: %v", ExchangeName, err)
		ch.Close()
		conn.Close()
		return
	}

	queue, err := ch.QueueDeclare(orderShippedQueue, true, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo declarar la cola %s: %v", orderShippedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	if err := ch.QueueBind(queue.Name, orderShippedRoutingKey, ExchangeName, false, nil); err != nil {
		log.Printf("[eventbus] no se pudo enlazar la cola %s: %v", orderShippedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("[eventbus] no se pudo iniciar el consumo de %s: %v", orderShippedQueue, err)
		ch.Close()
		conn.Close()
		return
	}

	log.Printf("[eventbus] escuchando %s en la cola %s", orderShippedRoutingKey, orderShippedQueue)

	go func() {
		for msg := range msgs {
			handleOrderShipped(pool, msg)
		}
	}()
}

func handleOrderShipped(pool *pgxpool.Pool, msg amqp.Delivery) {
	defer msg.Ack(false)

	var payload orderConfirmedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("[eventbus] mensaje order.shipped inválido: %v", err)
		return
	}

	data, err := json.Marshal(map[string]string{"asset_id": payload.AssetID, "order_id": payload.OrderID})
	if err != nil {
		log.Printf("[eventbus] no se pudo serializar data de la notificacion de envío: %v", err)
		return
	}

	const query = `
		INSERT INTO notifications (user_id, type, subtype, title, body, data)
		VALUES ($1, 'venta', 'pedido_enviado', $2, $3, $4)
	`
	const title = "Tu pedido fue enviado"
	const body = "El vendedor marcó tu pedido como enviado. Confírmalo cuando lo recibas."
	if _, err := pool.Exec(context.Background(), query, payload.BuyerID, title, body, data); err != nil {
		log.Printf("[eventbus] no se pudo crear la notificacion de envío para la orden %s: %v", payload.OrderID, err)
	}
}
