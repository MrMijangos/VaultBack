package eventbus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	nlpAnalyzedQueue      = "vault-community-service.nlp-analyzed"
	nlpAnalyzedRoutingKey = "nlp.analyzed"

	blockchainConfirmedQueue      = "vault-api.blockchain-confirmed"
	blockchainConfirmedRoutingKey = "blockchain.confirmed"
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
