package application

import (
	"context"
	"log"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/comments/domain/dto/request"
	"vault/src/features/comments/domain/dto/response"
	"vault/src/features/comments/domain/entities"
	"vault/src/features/comments/domain/repositories"
	notifentities "vault/src/features/notifications/domain/entities"
)

type CreateCommentUseCase struct {
	repo        repositories.CommentRepository
	moderation  *moderation.Client
	postAuthors repositories.PostAuthorProvider
	notifier    repositories.NotificationPublisher
}

func NewCreateCommentUseCase(
	repo repositories.CommentRepository,
	moderationClient *moderation.Client,
	postAuthors repositories.PostAuthorProvider,
	notifier repositories.NotificationPublisher,
) *CreateCommentUseCase {
	return &CreateCommentUseCase{repo: repo, moderation: moderationClient, postAuthors: postAuthors, notifier: notifier}
}

func (uc *CreateCommentUseCase) Execute(ctx context.Context, postID string, userID string, req request.CreateCommentRequest) (response.CommentResponse, error) {
	if err := req.Validate(); err != nil {
		return response.CommentResponse{}, err
	}

	commentID := uuid.NewString()

	result, err := uc.moderation.Analyze(ctx, commentID, "comment", req.Content)
	if err != nil {
		return response.CommentResponse{}, err
	}
	if result.IsToxic {
		return response.CommentResponse{}, moderation.ErrToxicContent
	}

	created, err := uc.repo.Create(ctx, entities.Comment{
		ID:            commentID,
		PostID:        postID,
		UserID:        userID,
		Content:       req.Content,
		ToxicityScore: &result.ToxicityScore,
		IsVisible:     true,
	})
	if err != nil {
		return response.CommentResponse{}, err
	}

	// Si la notificación falla, el comentario ya se guardó igual -- no vale
	// la pena tumbar la respuesta por esto, solo se registra.
	if authorID, err := uc.postAuthors.FindAuthorID(ctx, postID); err != nil {
		log.Printf("no se pudo obtener el autor de la publicacion %s: %v", postID, err)
	} else if authorID != userID {
		if _, err := uc.notifier.Create(ctx, notifentities.Notification{
			UserID:  authorID,
			Type:    "comunidad",
			Subtype: "comentario_nuevo",
			Title:   "Nuevo comentario",
			Body:    "Alguien comentó en tu publicación.",
		}); err != nil {
			log.Printf("no se pudo crear la notificacion de comentario nuevo: %v", err)
		}
	}

	return response.FromEntity(created), nil
}
