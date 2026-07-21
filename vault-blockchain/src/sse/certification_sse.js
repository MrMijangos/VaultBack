class CertificationSSE {
  constructor() {
    this.streams = new Map(); // asset_id -> [{ res, keepAlive }, ...]
  }

  addStream(assetId, res) {
    res.setHeader('Content-Type', 'text/event-stream');
    res.setHeader('Cache-Control', 'no-cache');
    res.setHeader('Connection', 'keep-alive');
    res.setHeader('X-Accel-Buffering', 'no');
    res.flushHeaders?.();

    res.write(
      `data: ${JSON.stringify({
        status: 'pending',
        message: 'Certificando en Vara Network...',
        asset_id: assetId,
      })}\n\n`
    );

    // Sin esto los proxies/load balancers cierran la conexión por inactividad.
    const keepAlive = setInterval(() => {
      res.write(': keepalive\n\n');
    }, 15000);

    if (!this.streams.has(assetId)) {
      this.streams.set(assetId, []);
    }
    this.streams.get(assetId).push({ res, keepAlive });

    res.on('close', () => {
      clearInterval(keepAlive);
      this.removeStream(assetId, res);
    });
  }

  notifyConfirmed(assetId, txId, confirmedAt) {
    const clients = this.streams.get(assetId) || [];
    const payload = JSON.stringify({
      status: 'confirmed',
      tx_id: txId,
      asset_id: assetId,
      confirmed_at: confirmedAt,
      network: process.env.VARA_NETWORK || 'testnet',
    });

    clients.forEach(({ res, keepAlive }) => {
      clearInterval(keepAlive);
      res.write(`event: confirmed\n`);
      res.write(`data: ${payload}\n\n`);
      res.end();
    });

    this.streams.delete(assetId);
  }

  notifyError(assetId, errorMsg) {
    const clients = this.streams.get(assetId) || [];
    clients.forEach(({ res, keepAlive }) => {
      clearInterval(keepAlive);
      res.write(`event: error\n`);
      res.write(
        `data: ${JSON.stringify({
          status: 'error',
          asset_id: assetId,
          message: errorMsg,
        })}\n\n`
      );
      res.end();
    });
    this.streams.delete(assetId);
  }

  removeStream(assetId, res) {
    const clients = this.streams.get(assetId) || [];
    this.streams.set(
      assetId,
      clients.filter((c) => c.res !== res)
    );
  }
}

module.exports = new CertificationSSE();
