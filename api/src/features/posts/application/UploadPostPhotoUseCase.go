package application

import (
	"context"
	"io"

	"vault/src/core/cloudinary"
	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/repositories"
)

type UploadPostPhotoUseCase struct {
	repo     repositories.PostRepository
	uploader *cloudinary.ImageUploader
}

func NewUploadPostPhotoUseCase(repo repositories.PostRepository, uploader *cloudinary.ImageUploader) *UploadPostPhotoUseCase {
	return &UploadPostPhotoUseCase{repo: repo, uploader: uploader}
}

func (uc *UploadPostPhotoUseCase) Execute(ctx context.Context, postID string, userID string, file io.Reader) (response.PostResponse, error) {
	post, err := uc.repo.FindByID(ctx, postID)
	if err != nil {
		return response.PostResponse{}, err
	}
	if post.UserID != userID {
		return response.PostResponse{}, repositories.ErrPostNotFound
	}

	url, err := uc.uploader.Upload(ctx, file, "vault/posts")
	if err != nil {
		return response.PostResponse{}, err
	}

	if _, err := uc.repo.AddPhoto(ctx, postID, url); err != nil {
		return response.PostResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByPostID(ctx, postID)
	if err != nil {
		return response.PostResponse{}, err
	}

	return response.FromEntity(post, photos), nil
}
