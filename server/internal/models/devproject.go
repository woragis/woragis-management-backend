package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Project struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string         `gorm:"size:200;not null" json:"name"`
	Slug             string         `gorm:"uniqueIndex;size:120;not null" json:"slug"`
	Description      string         `gorm:"type:text" json:"description"`
	ShortDescription string         `gorm:"column:short_description;type:text" json:"shortDescription"`
	LongDescription  string         `gorm:"column:long_description;type:text" json:"longDescription"`
	Status           string         `gorm:"size:32;not null;default:active" json:"status"`
	Stack            datatypes.JSON   `gorm:"type:jsonb" json:"stack"`
	RepoURL          string         `gorm:"column:repo_url;size:500" json:"repoUrl"`
	DemoURL          string         `gorm:"column:demo_url;size:500" json:"demoUrl"`
	GithubURL        string         `gorm:"column:github_url;size:500" json:"githubUrl"`
	Notes            string         `gorm:"type:text" json:"notes"`
	IsPublic         bool           `gorm:"column:is_public;not null;default:false" json:"isPublic"`
	Featured         bool           `gorm:"not null;default:false" json:"featured"`
	DisplayOrder     int            `gorm:"column:display_order;not null;default:0" json:"displayOrder"`
	PublicSlug       string         `gorm:"column:public_slug;size:120" json:"publicSlug"`
	CoverImageID     *uuid.UUID     `gorm:"column:cover_image_id;type:uuid" json:"coverImageId"`
	ParentProjectID  *uuid.UUID     `gorm:"column:parent_project_id;type:uuid" json:"parentProjectId"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	Links            []ProjectLink  `gorm:"foreignKey:ProjectID" json:"links,omitempty"`
	Domains          []ProjectDomain `gorm:"foreignKey:ProjectID" json:"domains,omitempty"`
	Gallery          []ProjectGallery `gorm:"foreignKey:ProjectID" json:"gallery,omitempty"`
	Envs             []ProjectEnv   `gorm:"foreignKey:ProjectID" json:"envs,omitempty"`
}

type ProjectLink struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID `gorm:"column:project_id;type:uuid;not null;index" json:"projectId"`
	Type        string    `gorm:"size:32;not null" json:"type"`
	URL         string    `gorm:"size:500;not null" json:"url"`
	Environment string    `gorm:"size:32;not null;default:production" json:"environment"`
	Label       string    `gorm:"size:200" json:"label"`
	IsPublic    bool      `gorm:"column:is_public;not null;default:false" json:"isPublic"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProjectDomain struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID uuid.UUID  `gorm:"column:project_id;type:uuid;not null;index" json:"projectId"`
	Domain    string     `gorm:"size:253;not null" json:"domain"`
	Registrar string     `gorm:"size:120" json:"registrar"`
	Purpose   string     `gorm:"size:64" json:"purpose"`
	ExpiresAt *time.Time `gorm:"column:expires_at" json:"expiresAt"`
	Notes     string     `gorm:"type:text" json:"notes"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type ProjectSecret struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID      uuid.UUID  `gorm:"column:project_id;type:uuid;not null;index" json:"projectId"`
	Name           string     `gorm:"size:200;not null" json:"name"`
	EncryptedValue string     `gorm:"column:encrypted_value;type:text;not null" json:"-"`
	Environment    string     `gorm:"size:32;not null;default:production" json:"environment"`
	Service        string     `gorm:"size:120" json:"service"`
	ExpiresAt      *time.Time `gorm:"column:expires_at" json:"expiresAt"`
	Notes          string     `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type ProjectSecretView struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"projectId"`
	Name        string     `json:"name"`
	Value       string     `json:"value,omitempty"`
	Environment string     `json:"environment"`
	Service     string     `json:"service"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type ProjectGallery struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID    uuid.UUID `gorm:"column:project_id;type:uuid;not null;index" json:"projectId"`
	MediaAssetID uuid.UUID `gorm:"column:media_asset_id;type:uuid;not null" json:"mediaAssetId"`
	DisplayOrder int       `gorm:"column:display_order;not null;default:0" json:"displayOrder"`
	Caption      string    `gorm:"size:300" json:"caption"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ProjectEnv struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID `gorm:"column:project_id;type:uuid;not null;index" json:"projectId"`
	Key         string    `gorm:"size:200;not null" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Environment string    `gorm:"size:32;not null;default:production" json:"environment"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
