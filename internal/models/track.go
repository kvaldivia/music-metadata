package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	SpotifyID string
	ISRC      string
	ImageURI  string
	Title     string
	ArtistID  uint
}
