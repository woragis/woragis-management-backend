package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaAsset struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Filename   string    `gorm:"size:255;not null" json:"filename"`
	MimeType   string    `gorm:"column:mime_type;size:120;not null" json:"mimeType"`
	SizeBytes  int64     `gorm:"column:size_bytes;not null" json:"sizeBytes"`
	StorageKey string    `gorm:"column:storage_key;size:500;not null" json:"-"`
	PublicURL  string    `gorm:"column:public_url;size:500;not null" json:"publicUrl"`
	AltText    string    `gorm:"column:alt_text;size:300" json:"altText"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	CreatedAt  time.Time `json:"createdAt"`
}
