package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	LeetcodeVideoStatusDraft     = "draft"
	LeetcodeVideoStatusPublished = "published"
	LeetcodeVideoStatusScheduled = "scheduled"

	ContentThumbnailStatusDraft      = "draft"
	ContentThumbnailStatusGenerating = "generating"
	ContentThumbnailStatusReady      = "ready"
	ContentThumbnailStatusApproved     = "approved"
	ContentThumbnailStatusFailed       = "failed"
)

type LeetcodeVideo struct {
	ID                    uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Title                 string          `gorm:"size:300;not null" json:"title"`
	Status                string          `gorm:"size:32;not null;default:draft" json:"status"`
	LeetcodeProblemNumber *int            `json:"leetcodeProblemNumber,omitempty"`
	LeetcodeSlug          *string         `gorm:"size:200" json:"leetcodeSlug,omitempty"`
	Difficulty            *string         `gorm:"size:32" json:"difficulty,omitempty"`
	Topics                json.RawMessage `gorm:"type:jsonb" json:"topics,omitempty"`
	Notes                 *string         `gorm:"type:text" json:"notes,omitempty"`
	YoutubeURL            *string         `gorm:"size:500" json:"youtubeUrl,omitempty"`
	PublishedAt           *time.Time      `json:"publishedAt,omitempty"`
	CreatedAt             time.Time       `json:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt"`

	Thumbnails []ContentThumbnail `gorm:"foreignKey:VideoID" json:"thumbnails,omitempty"`
}

type ContentThumbnail struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	VideoID           uuid.UUID       `gorm:"type:uuid;not null;index" json:"videoId"`
	Status            string          `gorm:"size:32;not null;default:draft" json:"status"`
	Prompt            string          `gorm:"type:text;not null" json:"prompt"`
	NegativePrompt    *string         `gorm:"type:text" json:"negativePrompt,omitempty"`
	Size              string          `gorm:"size:32;not null;default:1280x720" json:"size"`
	Quality           string          `gorm:"size:32;not null;default:high" json:"quality"`
	Model             string          `gorm:"size:64;not null;default:gpt-image-2" json:"model"`
	Mode              string          `gorm:"size:32;not null;default:edit" json:"mode"`
	ReferenceMediaIDs json.RawMessage `gorm:"type:jsonb" json:"referenceMediaIds,omitempty"`
	CreativesJobID    *uuid.UUID      `gorm:"type:uuid" json:"creativesJobId,omitempty"`
	OutputMediaID     *uuid.UUID      `gorm:"type:uuid" json:"outputMediaId,omitempty"`
	MetadataJSON      json.RawMessage `gorm:"type:jsonb" json:"metadata,omitempty"`
	ErrorMessage      *string         `gorm:"type:text" json:"errorMessage,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

type ContentPromptTemplate struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ChannelSlug    string    `gorm:"size:64;not null;index" json:"channelSlug"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	Slug           string    `gorm:"size:200;not null" json:"slug"`
	PromptTemplate string    `gorm:"type:text;not null" json:"promptTemplate"`
	IsDefault      bool      `gorm:"not null;default:false" json:"isDefault"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
