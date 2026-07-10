package repositories

import (
	"context"
	"errors"

	"vault/src/features/blockchaincertificates/domain/entities"
)

var ErrCertificateNotFound = errors.New("el certificado no existe")
var ErrTxIDAlreadyExists = errors.New("ese tx_id ya fue registrado")

type BlockchainCertificateRepository interface {
	Create(ctx context.Context, cert entities.BlockchainCertificate) (entities.BlockchainCertificate, error)
	ExistsByTxID(ctx context.Context, txID string) (bool, error)
	FindByAssetID(ctx context.Context, assetID string) ([]entities.BlockchainCertificate, error)
	FindByID(ctx context.Context, id string) (entities.BlockchainCertificate, error)
}
