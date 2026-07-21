const express = require('express');
const { sendCertificate } = require('../vara/smart_contract');
const certSSE = require('../sse/certification_sse');
const db = require('../db/supabase_client');
const publisher = require('../publishers/rabbitmq_publisher');

const router = express.Router();

// POST /certify -- llamado por vault-backend (Go) via HTTP interno.
router.post('/certify', (req, res) => {
  const { asset_id, owner_id, asset_hash, action } = req.body;

  if (!asset_id || !owner_id || !asset_hash || !action) {
    return res.status(400).json({ error: 'Faltan campos requeridos' });
  }

  // Responder de inmediato -- la certificación en Vara tarda 10-30s, no se
  // puede hacer esperar al caller.
  res.status(202).json({ message: 'Certificación iniciada', asset_id });

  certifyInBackground(asset_id, owner_id, asset_hash, action);
});

// GET /certify/status/:assetId -- Flutter se conecta aquí via SSE.
router.get('/certify/status/:assetId', (req, res) => {
  certSSE.addStream(req.params.assetId, res);
});

async function certifyInBackground(assetId, ownerId, assetHash, action) {
  try {
    const txId = await sendCertificate({ assetId, ownerId, assetHash, action });
    const confirmedAt = new Date().toISOString();
    const network = process.env.VARA_NETWORK || 'testnet';

    await db.insertCertificate({ assetId, ownerId, txId, assetHash, action, network, confirmedAt });
    await db.updateAssetBlockchainInfo({ assetId, txId, assetHash });

    certSSE.notifyConfirmed(assetId, txId, confirmedAt);

    await publisher.publish('blockchain.confirmed', {
      asset_id: assetId,
      tx_id: txId,
      confirmed_at: confirmedAt,
      action,
    });

    console.log(`[Vara] Asset ${assetId} certificado: ${txId}`);
  } catch (error) {
    console.error(`[Vara] Error certificando ${assetId}:`, error.message);
    certSSE.notifyError(assetId, error.message);
  }
}

module.exports = { router, certifyInBackground };
