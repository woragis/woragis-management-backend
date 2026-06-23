package models

import (
	"time"

	"github.com/google/uuid"
)

// DefaultPersonalityID is the singleton row for agent personality (v1).
var DefaultPersonalityID = uuid.MustParse("00000000-0000-4000-8000-000000000001")

type AgentPersonality struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AssistantName     string     `gorm:"column:assistant_name;size:120;not null;default:Assistente" json:"assistantName"`
	AssistantNameSetAt *time.Time `gorm:"column:assistant_name_set_at" json:"assistantNameSetAt"`
	GreetingMorning   string     `gorm:"column:greeting_morning;size:300;not null;default:Bom dia" json:"greetingMorning"`
	GreetingAfternoon string     `gorm:"column:greeting_afternoon;size:300;not null;default:Boa tarde" json:"greetingAfternoon"`
	GreetingEvening   string     `gorm:"column:greeting_evening;size:300;not null;default:Boa noite" json:"greetingEvening"`
	GreetingEnabled   bool       `gorm:"column:greeting_enabled;not null;default:false" json:"greetingEnabled"`
	SystemPromptExtra string     `gorm:"column:system_prompt_extra;type:text" json:"systemPromptExtra"`
	VoiceID           string     `gorm:"column:voice_id;size:64;not null;default:alloy" json:"voiceId"`
	Language          string     `gorm:"size:16;not null;default:pt-BR" json:"language"`
	Timezone          string     `gorm:"size:64;not null;default:America/Sao_Paulo" json:"timezone"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

func DefaultAgentPersonality() AgentPersonality {
	return AgentPersonality{
		ID:                DefaultPersonalityID,
		AssistantName:     "Assistente",
		GreetingMorning:   "Bom dia",
		GreetingAfternoon: "Boa tarde",
		GreetingEvening:   "Boa noite",
		GreetingEnabled:   false,
		VoiceID:           "alloy",
		Language:          "pt-BR",
		Timezone:          "America/Sao_Paulo",
	}
}
