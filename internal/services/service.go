package services

import (
	"context"

	"github.com/kvaldivia/music-metadata/internal/models"
)

type Service interface {
	Get(ctx context.Context, id string) (*models.Track, *models.Artist, error)
	Search(ctx context.Context, q string) ([]*models.Track, error)
}
