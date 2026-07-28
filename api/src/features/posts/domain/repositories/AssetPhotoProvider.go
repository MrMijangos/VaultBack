package repositories

import (
	"context"

	assetentities "vault/src/features/assets/domain/entities"
)

// AssetPhotoProvider es el puerto hacia el feature de assets -- mismo
// patrón que NotificationPublisher en chat/: en infrastructure/dependencies.go
// se pasa la implementación real de assets/infrastructure/adapters
// directamente, que ya lo satisface por estructura sin que posts dependa de
// su application/infraestructura. Se usa para copiar las fotos del activo
// al post cuando se crea con asset_id (ver "Publicar en el Feed" en
// vault-app, que hasta ahora solo mandaba el asset_id sin ninguna foto).
type AssetPhotoProvider interface {
	FindPhotosByAssetID(ctx context.Context, assetID string) ([]assetentities.AssetPhoto, error)
}
