const express = require('express');
const db = require('../db/supabase_client');

const router = express.Router();

// GET /verify/:assetId -- historial de certificados de un asset.
router.get('/verify/:assetId', async (req, res) => {
  try {
    const certificates = await db.getCertificatesByAsset(req.params.assetId);
    return res.json({
      asset_id: req.params.assetId,
      total: certificates.length,
      is_certified: certificates.length > 0,
      certificates,
    });
  } catch (error) {
    return res.status(500).json({ error: error.message });
  }
});

module.exports = router;
