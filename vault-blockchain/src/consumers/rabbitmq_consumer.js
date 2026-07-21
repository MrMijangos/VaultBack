// rabbitmq_consumer.js -- se conecta al mismo exchange topic "vault.events"
// que usan api/ y payment/ (no colas sueltas por fuera del exchange).
//
// IMPORTANTE: al momento de escribir esto, api/ TODAVÍA NO publica ni
// "asset.created" ni "maintenance.registered" -- solo publica "asset.updated"
// (para el servicio de ML, con un payload distinto: {event_type, user_id,
// source_id}) y no publica nada al crear un maintenance log. Este consumer
// se queda escuchando sin recibir nada hasta que se agregue esa publicación
// en api/ (CreateAssetUseCase.go / CreateMaintenanceLogUseCase.go), usando
// el payload que se documenta abajo en handleAssetCreated/handleMaintenanceRegistered.
const amqp = require('amqplib');
const crypto = require('crypto');

const EXCHANGE_NAME = 'vault.events';

// generateAssetHash -- payload debe traer los mismos campos que expone
// AssetResponse en api/ (id, user_id, name, category, created_at), no los
// nombres inventados de una spec genérica.
function generateAssetHash(asset) {
  const payload = JSON.stringify({
    asset_id: asset.asset_id,
    user_id: asset.user_id,
    name: asset.name,
    category: asset.category,
    created_at: asset.created_at,
  });
  return crypto.createHash('sha256').update(payload).digest('hex');
}

function generateMaintenanceHash(maintenance) {
  const payload = JSON.stringify({
    maintenance_id: maintenance.maintenance_id,
    asset_id: maintenance.asset_id,
  });
  return crypto.createHash('sha256').update(payload).digest('hex');
}

async function bindQueue(channel, queueName, routingKey, handler) {
  await channel.assertQueue(queueName, { durable: true });
  await channel.bindQueue(queueName, EXCHANGE_NAME, routingKey);

  channel.consume(queueName, async (msg) => {
    if (!msg) return;
    try {
      const data = JSON.parse(msg.content.toString());
      await handler(data);
      channel.ack(msg);
    } catch (err) {
      console.error(`[Consumer] Error en ${routingKey}:`, err.message);
      channel.nack(msg, false, true); // requeue
    }
  });
}

async function startConsumer(certifyFn) {
  if (!process.env.RABBITMQ_URL) {
    console.log('[RabbitMQ] RABBITMQ_URL no configurado -- consumer deshabilitado');
    return;
  }

  const conn = await amqp.connect(process.env.RABBITMQ_URL);
  const channel = await conn.createChannel();
  await channel.assertExchange(EXCHANGE_NAME, 'topic', { durable: true });

  await bindQueue(channel, 'vault-blockchain.asset-created', 'asset.created', async (data) => {
    const hash = generateAssetHash(data);
    await certifyFn(data.asset_id, data.user_id, hash, 'REGISTERED');
  });

  await bindQueue(
    channel,
    'vault-blockchain.maintenance-registered',
    'maintenance.registered',
    async (data) => {
      const action = data.type === 'restauracion' ? 'RESTORED' : 'MAINTAINED';
      const hash = generateMaintenanceHash(data);
      await certifyFn(data.asset_id, data.user_id, hash, action);
    }
  );

  conn.on('close', () => console.log('[RabbitMQ] Consumer desconectado'));

  console.log('[RabbitMQ] Consumer iniciado -- escuchando asset.created y maintenance.registered');
}

module.exports = { startConsumer };
