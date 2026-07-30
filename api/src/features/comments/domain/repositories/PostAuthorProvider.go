package repositories

import "context"

// PostAuthorProvider es el puerto hacia posts/ -- mismo patrón que
// AssetPhotoProvider en posts/domain/repositories: en
// infrastructure/dependencies.go se pasa PostgreSQLPostRepository
// directamente, que lo satisface por estructura sin que comments dependa de
// la application/infraestructura de posts.
type PostAuthorProvider interface {
	FindAuthorID(ctx context.Context, postID string) (string, error)
}
