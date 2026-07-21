// supabase_client.js -- pese al nombre (por consistencia con la carpeta
// db/), esto conecta directo a Postgres via DATABASE_URL con "pg", igual
// que api/ y payment/ hacen con sus propias variables DB_HOST/DB_USER/etc.
// No se usa el SDK REST de @supabase/supabase-js porque lo único que nos
// dieron fue el connection string del pooler, no una SUPABASE_URL/API key.
const { Pool } = require('pg');

// El pooler de Supabase usa un certificado que Node no valida por default.
// Ojo: si el connection string trae "?sslmode=require", pg-connection-string
// lo trata como alias de "verify-full" e IGNORA el `ssl: {...}` de abajo --
// hay que quitar ese query param para que rejectUnauthorized:false aplique.
const connectionString = (process.env.DATABASE_URL || '').replace(/[?&]sslmode=[^&]*/, '');
const pool = new Pool({
  connectionString,
  ssl: { rejectUnauthorized: false },
});

async function insertCertificate({ assetId, ownerId, txId, assetHash, action, network, confirmedAt }) {
  await pool.query(
    `INSERT INTO blockchain_certificates (asset_id, owner_id, tx_id, asset_hash, action, network, confirmed_at)
     VALUES ($1, $2, $3, $4, $5, $6, $7)`,
    [assetId, ownerId, txId, assetHash, action, network, confirmedAt]
  );
}

async function updateAssetBlockchainInfo({ assetId, txId, assetHash }) {
  await pool.query(`UPDATE assets SET blockchain_tx_id = $1, blockchain_hash = $2 WHERE id = $3`, [
    txId,
    assetHash,
    assetId,
  ]);
}

async function getCertificatesByAsset(assetId) {
  const { rows } = await pool.query(
    `SELECT asset_id, owner_id, tx_id, asset_hash, action, network, confirmed_at
     FROM blockchain_certificates
     WHERE asset_id = $1
     ORDER BY confirmed_at ASC`,
    [assetId]
  );
  return rows;
}

module.exports = { pool, insertCertificate, updateAssetBlockchainInfo, getCertificatesByAsset };
