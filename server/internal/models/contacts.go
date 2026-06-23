package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Contact struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string         `gorm:"size:200;not null;index" json:"name"`
	DisplayName     string         `gorm:"column:display_name;size:300" json:"displayName"`
	Email           string         `gorm:"size:320;index" json:"email"`
	Phone           string         `gorm:"size:64;index" json:"phone"`
	Telegram        string         `gorm:"size:128" json:"telegram"`
	Whatsapp        string         `gorm:"size:64" json:"whatsapp"`
	Organization    string         `gorm:"size:200;index" json:"organization"`
	RoleTitle       string         `gorm:"column:role_title;size:200" json:"roleTitle"`
	Relationship    string         `gorm:"size:32;not null;default:other;index" json:"relationship"`
	Stage           string         `gorm:"size:32;not null;default:cold;index" json:"stage"`
	Source          string         `gorm:"size:200" json:"source"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Tags            datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	ProjectID       *uuid.UUID     `gorm:"column:project_id;type:uuid;index" json:"projectId"`
	LastContactedAt *time.Time     `gorm:"column:last_contacted_at" json:"lastContactedAt"`
	NextFollowUpAt  *time.Time     `gorm:"column:next_follow_up_at" json:"nextFollowUpAt"`
	Active          bool           `gorm:"not null;default:true;index" json:"active"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

type ContactInteraction struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ContactID  uuid.UUID `gorm:"column:contact_id;type:uuid;not null;index" json:"contactId"`
	Type       string    `gorm:"size:32;not null" json:"type"`
	Channel    string    `gorm:"size:32;not null;default:other" json:"channel"`
	Summary    string    `gorm:"type:text" json:"summary"`
	HappenedAt time.Time `gorm:"column:happened_at;not null;index" json:"happenedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}
