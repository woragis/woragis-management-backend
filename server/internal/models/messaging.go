package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	ChannelWhatsApp = "whatsapp"
	ChannelTelegram = "telegram"

	ComposeModeStatic     = "static"
	ComposeModeAIAssisted = "ai_assisted"
)

type ChannelDestination struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Channel          string         `gorm:"size:32;not null;index" json:"channel"`
	ExternalID       string         `gorm:"column:external_id;size:256;not null;index" json:"externalId"`
	Name             string         `gorm:"size:200;not null" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	Responsibilities string         `gorm:"type:text" json:"responsibilities"`
	Tags             datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	Active           bool           `gorm:"not null;default:true;index" json:"active"`
	Metadata         datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

type MessageTemplate struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	DestinationID *uuid.UUID `gorm:"column:destination_id;type:uuid;index" json:"destinationId"`
	ProgramSlug   string     `gorm:"column:program_slug;size:64;index" json:"programSlug"`
	Slug          string     `gorm:"size:64;not null;index:idx_msg_tpl_program_slug" json:"slug"`
	Name          string     `gorm:"size:200;not null" json:"name"`
	Body          string     `gorm:"type:text;not null" json:"body"`
	ComposeMode   string     `gorm:"column:compose_mode;size:32;not null;default:static" json:"composeMode"`
	AIPromptHint  string     `gorm:"column:ai_prompt_hint;type:text" json:"aiPromptHint"`
	Active        bool       `gorm:"not null;default:true" json:"active"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type ScheduledJob struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string     `gorm:"size:200;not null" json:"name"`
	DestinationID uuid.UUID  `gorm:"column:destination_id;type:uuid;not null;index" json:"destinationId"`
	TemplateSlug  string     `gorm:"column:template_slug;size:64" json:"templateSlug"`
	ProgramAction string     `gorm:"column:program_action;size:64;index" json:"programAction"`
	CronExpr      string     `gorm:"column:cron_expr;size:64;not null" json:"cronExpr"`
	Timezone      string     `gorm:"size:64;not null;default:America/Sao_Paulo" json:"timezone"`
	Enabled       bool       `gorm:"not null;default:true;index" json:"enabled"`
	LastRunAt     *time.Time `gorm:"column:last_run_at" json:"lastRunAt"`
	NextRunAt     *time.Time `gorm:"column:next_run_at;index" json:"nextRunAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type MessageDelivery struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	JobID         *uuid.UUID `gorm:"column:job_id;type:uuid;index" json:"jobId"`
	DestinationID uuid.UUID  `gorm:"column:destination_id;type:uuid;not null;index" json:"destinationId"`
	Channel       string     `gorm:"size:32;not null" json:"channel"`
	ExternalID    string     `gorm:"column:external_id;size:256;not null" json:"externalId"`
	TemplateSlug  string     `gorm:"column:template_slug;size:64" json:"templateSlug"`
	Body          string     `gorm:"type:text;not null" json:"body"`
	Status        string     `gorm:"size:16;not null;default:sent;index" json:"status"`
	ErrorMessage  string     `gorm:"column:error_message;type:text" json:"errorMessage"`
	ExternalRef   string     `gorm:"column:external_ref;size:128" json:"externalRef"`
	SentAt        time.Time  `gorm:"column:sent_at;not null;index" json:"sentAt"`
	CreatedAt     time.Time  `json:"createdAt"`
}
