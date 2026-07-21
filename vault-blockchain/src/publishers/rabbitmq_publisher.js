// rabbitmq_publisher.js -- publica en el mismo exchange topic "vault.events"
// que usan api/ y payment/ (ver api/src/core/eventbus/Publisher.go), no crea
// uno nuevo. La routing key es el propio nombre del evento, igual que allá.
const amqp = require('amqplib');

const EXCHANGE_NAME = 'vault.events';

let channel = null;

async function getChannel() {
  if (channel) return channel;
  if (!process.env.RABBITMQ_URL) return null;

  const conn = await amqp.connect(process.env.RABBITMQ_URL);
  channel = await conn.createChannel();
  await channel.assertExchange(EXCHANGE_NAME, 'topic', { durable: true });

  conn.on('close', () => {
    console.log('[RabbitMQ] Publisher desconectado');
    channel = null;
  });

  return channel;
}

// publish -- si RabbitMQ no está configurado o falla, solo loguea (mismo
// comportamiento "noop" que api/ y payment/: el servicio sigue funcionando).
async function publish(eventType, payload) {
  try {
    const ch = await getChannel();
    if (!ch) {
      console.log(`[RabbitMQ] No configurado, no se publica ${eventType}`);
      return;
    }
    const body = Buffer.from(JSON.stringify({ event_type: eventType, ...payload }));
    ch.publish(EXCHANGE_NAME, eventType, body, { contentType: 'application/json' });
  } catch (err) {
    console.error(`[RabbitMQ] Error publicando ${eventType}:`, err.message);
  }
}

module.exports = { publish };
