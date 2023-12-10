package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model
	Name  string  `json:"name,omitempty"`
	Songs []Track `json:"songs" gorm:"foreignKey:ArtistRef"`
}
