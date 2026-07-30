package application

import (
	"context"
	"log"

	notifentities "vault/src/features/notifications/domain/entities"
	"vault/src/features/posts/domain/repositories"
)

type LikePostUseCase struct {
	repo     repositories.PostRepository
	notifier repositories.NotificationPublisher
}

func NewLikePostUseCase(repo repositories.PostRepository, notifier repositories.NotificationPublisher) *LikePostUseCase {
	return &LikePostUseCase{repo: repo, notifier: notifier}
}

func (uc *LikePostUseCase) Execute(ctx context.Context, postID string, userID string) error {
	if err := uc.repo.Like(ctx, postID, userID); err != nil {
		return err
	}

	// Si la notificación falla, el like ya quedó registrado igual -- no vale
	// la pena tumbar la respuesta por esto, solo se registra.
	post, err := uc.repo.FindByID(ctx, postID)
	if err != nil {
		log.Printf("no se pudo obtener el autor de la publicacion %s: %v", postID, err)
		return nil
	}
	if post.UserID == userID {
		return nil
	}
	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  post.UserID,
		Type:    "comunidad",
		Subtype: "likes_post",
		Title:   "Nuevo like",
		Body:    "A alguien le gustó tu publicación.",
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de like nuevo: %v", err)
	}
	return nil
}
