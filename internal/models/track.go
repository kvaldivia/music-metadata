package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	SpotifyID string `json:"spotify_id"`
	ISRC      string `json:"string" `
	ImageURI  string `json:"imageUri,omitempty"`
	Title     string `json:"title,omitempty"`
	ArtistRef string
}
