const { CreateType } = require('@gear-js/api');
const { getApi, getAccount } = require('./gear_api');

// Debe coincidir exactamente con CertificateMessage en
// contracts/asset_certificate/src/lib.rs -- Gear decodifica el payload como
// bytes SCALE-encoded, no como JSON, así que aquí lo codificamos con el
// mismo shape que espera el contrato.
const CUSTOM_TYPES = {
  CertificateMessage: {
    asset_id: 'String',
    owner_id: 'String',
    asset_hash: 'String',
    action: 'String',
  },
};

const createType = new CreateType(CUSTOM_TYPES);

// sendCertificate -- envía un CertificateMessage al contrato desplegado en
// VARA_CONTRACT_ID y espera a que el bloque se finalice. Devuelve el hash
// del bloque como tx_id.
async function sendCertificate({ assetId, ownerId, assetHash, action }) {
  const contractId = process.env.VARA_CONTRACT_ID;
  if (!contractId) {
    throw new Error('VARA_CONTRACT_ID no está configurado');
  }

  const api = await getApi();
  const account = getAccount();

  const payloadHex = createType
    .create('CertificateMessage', {
      asset_id: assetId,
      owner_id: ownerId,
      asset_hash: assetHash,
      action,
    })
    .toHex();

  const gasLimit = await resolveGasLimit(api, account, contractId, payloadHex);

  const extrinsic = await api.message.send({
    destination: contractId,
    payload: payloadHex,
    gasLimit,
    value: 0,
  });

  return new Promise((resolve, reject) => {
    extrinsic
      .signAndSend(account, ({ status, events }) => {
        if (status.isInvalid || status.isDropped || status.isRetracted) {
          reject(new Error(`Transacción rechazada por Vara Network (${status.type})`));
          return;
        }

        const failed = (events || []).find(
          ({ event }) => event.section === 'system' && event.method === 'ExtrinsicFailed'
        );
        if (failed) {
          reject(new Error('El contrato rechazó el mensaje (ExtrinsicFailed)'));
          return;
        }

        if (status.isFinalized) {
          resolve(status.asFinalized.toHex());
        }
      })
      .catch(reject);
  });
}

// resolveGasLimit -- usa VARA_GAS_LIMIT si está fijado en .env, si no estima
// el gas real necesario con calculateGas.handle (+20% de margen).
async function resolveGasLimit(api, account, contractId, payloadHex) {
  const fixed = process.env.VARA_GAS_LIMIT;
  if (fixed) return BigInt(fixed);

  const gasInfo = await api.program.calculateGas.handle(
    account.address,
    contractId,
    payloadHex,
    0,
    false
  );
  return (gasInfo.min_limit.toBigInt() * 120n) / 100n;
}

module.exports = { sendCertificate };
