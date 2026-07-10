package cloudinary

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"vault/src/core/config"
)

type ImageUploader struct {
	client *cloudinary.Cloudinary
}

func NewImageUploader(cfg *config.Config) (*ImageUploader, error) {
	client, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		return nil, fmt.Errorf("no se pudo inicializar cloudinary: %w", err)
	}
	return &ImageUploader{client: client}, nil
}

func (u *ImageUploader) Upload(ctx context.Context, file io.Reader, folder string) (string, error) {
	result, err := u.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         folder,
		ResourceType:   "image",
		UniqueFilename: boolPtr(true),
		Overwrite:      boolPtr(false),
	})
	if err != nil {
		return "", fmt.Errorf("no se pudo subir la imagen: %w", err)
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf("cloudinary rechazo la imagen: %s", result.Error.Message)
	}

	return result.SecureURL, nil
}

func boolPtr(b bool) *bool {
	return &b
}
