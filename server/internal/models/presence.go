package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SocialPlatformLinkedIn = "linkedin"
	SocialPlatformReddit   = "reddit"
	SocialPlatformTwitter  = "twitter"

	SocialGoalJobHunting    = "job_hunting"
	SocialGoalRevenue       = "revenue"
	SocialGoalLaunch        = "launch"
	SocialGoalVisibility    = "visibility"
	SocialGoalAcademic      = "academic"
	SocialGoalCommunity     = "community"

	SocialPostStatusDraft     = "draft"
	SocialPostStatusScheduled = "scheduled"
	SocialPostStatusPublished = "published"
	SocialPostStatusCancelled = "cancelled"
)

type SocialCampaign struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string     `gorm:"size:200;not null" json:"name"`
	Goal        string     `gorm:"size:32;not null;index" json:"goal"`
	Description string     `gorm:"type:text" json:"description"`
	ProjectID   *uuid.UUID `gorm:"column:project_id;type:uuid;index" json:"projectId"`
	StartDate   *time.Time `gorm:"column:start_date;type:date" json:"startDate"`
	EndDate     *time.Time `gorm:"column:end_date;type:date" json:"endDate"`
	Active      bool       `gorm:"not null;default:true;index" json:"active"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type PostTemplate struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Slug      string    `gorm:"size:64;not null;uniqueIndex" json:"slug"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	Platform  string    `gorm:"size:32;not null;default:any;index" json:"platform"`
	Goal      string    `gorm:"size:32;not null;default:visibility;index" json:"goal"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	Active    bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SocialPost struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID    *uuid.UUID `gorm:"column:project_id;type:uuid;index" json:"projectId"`
	CampaignID   *uuid.UUID `gorm:"column:campaign_id;type:uuid;index" json:"campaignId"`
	Platform     string     `gorm:"size:32;not null;index" json:"platform"`
	Goal         string     `gorm:"size:32;not null;index" json:"goal"`
	Status       string     `gorm:"size:32;not null;default:draft;index" json:"status"`
	Title        string     `gorm:"size:300" json:"title"`
	Body         string     `gorm:"type:text;not null" json:"body"`
	Hook         string     `gorm:"size:500" json:"hook"`
	CTA          string     `gorm:"column:cta;size:500" json:"cta"`
	TemplateSlug string     `gorm:"column:template_slug;size:64" json:"templateSlug"`
	ScheduledAt  *time.Time `gorm:"column:scheduled_at;index" json:"scheduledAt"`
	PublishedAt  *time.Time `gorm:"column:published_at" json:"publishedAt"`
	PublishedURL string     `gorm:"column:published_url;size:500" json:"publishedUrl"`
	Notes        string     `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
