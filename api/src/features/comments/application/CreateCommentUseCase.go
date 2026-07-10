package application

import (
	"context"

	"vault/src/features/comments/domain/dto/request"
	"vault/src/features/comments/domain/dto/response"
	"vault/src/features/comments/domain/entities"
	"vault/src/features/comments/domain/repositories"
)

type CreateCommentUseCase struct {
	repo repositories.CommentRepository
}

func NewCreateCommentUseCase(repo repositories.CommentRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{repo: repo}
}

func (uc *CreateCommentUseCase) Execute(ctx context.Context, postID string, userID string, req request.CreateCommentRequest) (response.CommentResponse, error) {
	if err := req.Validate(); err != nil {
		return response.CommentResponse{}, err
	}

	created, err := uc.repo.Create(ctx, entities.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	})
	if err != nil {
		return response.CommentResponse{}, err
	}

	return response.FromEntity(created), nil
}
