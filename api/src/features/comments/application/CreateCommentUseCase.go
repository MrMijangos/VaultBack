package application

import (
	"context"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/comments/domain/dto/request"
	"vault/src/features/comments/domain/dto/response"
	"vault/src/features/comments/domain/entities"
	"vault/src/features/comments/domain/repositories"
)

type CreateCommentUseCase struct {
	repo       repositories.CommentRepository
	moderation *moderation.Client
}

func NewCreateCommentUseCase(repo repositories.CommentRepository, moderationClient *moderation.Client) *CreateCommentUseCase {
	return &CreateCommentUseCase{repo: repo, moderation: moderationClient}
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

	return response.FromEntity(created), nil
}
