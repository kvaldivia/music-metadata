package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	ISRC     string   `json:"string"`
	ImageURI string   `json:"imageUri,omitempty"`
	Title    string   `json:"title,omitempty"`
	Artists  []string `json:"artistNameList,omitempty"`
}
