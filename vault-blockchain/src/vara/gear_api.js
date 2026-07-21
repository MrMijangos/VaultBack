const { GearApi, GearKeyring } = require('@gear-js/api');

let api = null;
let account = null;

async function getApi() {
  if (api && api.isConnected) return api;

  api = await GearApi.create({
    providerAddress: process.env.VARA_NODE_URL,
  });

  api.provider.on('disconnected', () => {
    console.log('[Vara] Desconectado -- se reconectará en la próxima solicitud');
    api = null;
  });

  console.log(`[Vara] Conectado a ${process.env.VARA_NODE_URL}`);
  return api;
}

function getAccount() {
  if (account) return account;
  if (!process.env.VARA_MNEMONIC) {
    throw new Error('VARA_MNEMONIC no está configurado');
  }
  account = GearKeyring.fromMnemonic(process.env.VARA_MNEMONIC, 'vault-blockchain');
  return account;
}

module.exports = { getApi, getAccount };
