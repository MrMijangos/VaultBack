#![no_std]

use gstd::{exec, msg, prelude::*};
use parity_scale_codec::{Decode, Encode};
use scale_info::TypeInfo;

// El orden de los campos importa: el codec SCALE es posicional, así que
// debe coincidir exactamente con CUSTOM_TYPES.CertificateMessage en
// vault-blockchain/src/vara/smart_contract.js (asset_id, owner_id,
// asset_hash, action, en ese orden).
#[derive(Encode, Decode, TypeInfo, Debug, Clone)]
pub struct CertificateMessage {
    pub asset_id: String,
    pub owner_id: String,
    pub asset_hash: String,
    pub action: String,
}

#[derive(Encode, Decode, TypeInfo, Debug, Clone)]
pub struct AssetCertificate {
    pub asset_id: String,
    pub owner_id: String,
    pub asset_hash: String,
    pub action: String,
    pub timestamp: u64,
}

static mut CERTIFICATES: Vec<AssetCertificate> = Vec::new();

#[unsafe(no_mangle)]
extern "C" fn handle() {
    let message: CertificateMessage = msg::load().expect("No se pudo decodificar CertificateMessage");

    let certificate = AssetCertificate {
        asset_id: message.asset_id,
        owner_id: message.owner_id,
        asset_hash: message.asset_hash,
        action: message.action,
        timestamp: exec::block_timestamp(),
    };

    unsafe {
        CERTIFICATES.push(certificate);
    }

    msg::reply(String::from("OK"), 0).expect("No se pudo responder");
}

#[unsafe(no_mangle)]
extern "C" fn state() {
    let certificates = unsafe { CERTIFICATES.clone() };
    msg::reply(certificates, 0).expect("No se pudo devolver el state");
}
