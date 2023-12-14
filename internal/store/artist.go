package store

import (
	"context"

	"github.com/kvaldivia/music-metadata/internal/models"
	"gorm.io/gorm"
)

type Artist interface {
	Record(ctx context.Context, artists []*models.Artist) error
	Save(ctx context.Context, artist *models.Artist) error
	Create(ctx context.Context, o interface{}) error
	Find(ctx context.Context, spotifyId string) (*models.Artist, error)
	All(ctx context.Context) ([]*models.Artist, error)
}

var _ Artist = &artist{}

func NewArtistStore(db *gorm.DB) Artist {
	return &artist{db: db}
}

type artist struct {
	db *gorm.DB
}

func (t *artist) Record(ctx context.Context, artists []*models.Artist) error {
	return t.db.WithContext(ctx).Save(artists).Error
}

func (t *artist) Save(ctx context.Context, artist *models.Artist) error {
	return t.db.WithContext(ctx).Save(artist).Error
}

func (t *artist) Create(ctx context.Context, o interface{}) error {
	return t.db.WithContext(ctx).Create(o).Error
}

func (t *artist) Find(ctx context.Context, spotifyId string) (*models.Artist, error) {
	var matchedArtist models.Artist
	err := t.db.WithContext(ctx).First(&matchedArtist, "spotify_id = ?", spotifyId).Error
	return &matchedArtist, err
}

// TODO: optimize with pagination
func (t *artist) All(ctx context.Context) ([]*models.Artist, error) {
	var artists []*models.Artist
	err := t.db.WithContext(ctx).Find(artists).Error

	return artists, err
}
