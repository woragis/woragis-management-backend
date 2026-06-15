package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Profile struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Slug           string         `gorm:"uniqueIndex;size:64;not null;default:default" json:"slug"`
	DisplayName    string         `gorm:"column:display_name;size:200;not null" json:"displayName"`
	Headline       string         `gorm:"size:300" json:"headline"`
	Bio            string         `gorm:"type:text" json:"bio"`
	AvatarID       *uuid.UUID     `gorm:"column:avatar_id;type:uuid" json:"avatarId"`
	Location       string         `gorm:"size:200" json:"location"`
	Availability   string         `gorm:"size:64;not null;default:not_available" json:"availability"`
	ResumeAssetID  *uuid.UUID     `gorm:"column:resume_asset_id;type:uuid" json:"resumeAssetId"`
	SocialLinks    datatypes.JSON `gorm:"column:social_links;type:jsonb" json:"socialLinks"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type SocialLink struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
