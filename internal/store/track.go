package store

import (
	"context"

	"github.com/kvaldivia/music-metadata/internal/models"
	"gorm.io/gorm"
)

type Track interface {
	Record(ctx context.Context, tracks []*models.Track) error
	Save(ctx context.Context, track *models.Track) error
	Find(ctx context.Context, isrc string) (*models.Track, error)
	AllByArtist(ctx context.Context, artistName string) ([]*models.Track, error)
	All(ctx context.Context) ([]*models.Track, error)
}

var _ Track = &track{}

func NewTrackStore(db *gorm.DB) Track {
	return &track{db: db}
}

type track struct {
	db *gorm.DB
}

func (t *track) Record(ctx context.Context, tracks []*models.Track) error {
	return t.db.WithContext(ctx).Save(tracks).Error
}

func (t *track) Save(ctx context.Context, track *models.Track) error {
	return t.db.WithContext(ctx).Save(track).Error
}

func (t *track) Find(ctx context.Context, isrc string) (*models.Track, error) {
	var matchedTrack models.Track
	err := t.db.WithContext(ctx).First(&matchedTrack, "isrc = ?", isrc).Error
	return &matchedTrack, err
}

// TODO(kvaldivia): optimize with pagination
func (t *track) AllByArtist(ctx context.Context, artistName string) ([]*models.Track, error) {
	var tracks []*models.Track
	var artist *models.Artist
	var err error

	err = t.db.WithContext(ctx).First(artist, "name=?", artistName).Error
	if err != nil {
		return nil, err
	}

	err = t.db.WithContext(ctx).InnerJoins("Artists").Find(tracks, "artist_ref=?", artist.ID).Error

	return tracks, err
}

// TODO: optimize with pagination
func (t *track) All(ctx context.Context) ([]*models.Track, error) {
	var tracks []*models.Track
	err := t.db.WithContext(ctx).Find(tracks).Error

	return tracks, err
}
