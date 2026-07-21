require('dotenv').config();
const express = require('express');
const { router: certifyRouter, certifyInBackground } = require('./routes/certify_routes');
const verifyRoutes = require('./routes/verify_routes');
const { startConsumer } = require('./consumers/rabbitmq_consumer');

const app = express();
app.use(express.json());

app.use('/', certifyRouter);
app.use('/', verifyRoutes);
app.get('/health', (_, res) => res.json({ status: 'ok', service: 'vault-blockchain' }));

async function main() {
  const PORT = process.env.PORT || 3001;
  app.listen(PORT, () => {
    console.log(`[vault-blockchain] Corriendo en puerto ${PORT}`);
  });

  // Si RabbitMQ no está disponible el servicio sigue funcionando (certify
  // via HTTP directo y verify siguen respondiendo), solo no hay consumer.
  try {
    await startConsumer(certifyInBackground);
  } catch (err) {
    console.error('[RabbitMQ] No se pudo iniciar el consumer:', err.message);
  }
}

main().catch((err) => {
  console.error('[vault-blockchain] Error fatal al iniciar:', err);
  process.exit(1);
});
