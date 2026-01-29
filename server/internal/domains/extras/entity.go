package extras

import (
    "time"

    "github.com/google/uuid"
)

// Extra represents a small user-provided note or miscellaneous item.
type Extra struct {
    ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
    Category  string    `gorm:"column:category;size:128" json:"category,omitempty"`
    Text      string    `gorm:"column:text;type:text" json:"text,omitempty"`
    Ordinal   int       `gorm:"column:ordinal;default:0" json:"ordinal,omitempty"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName sets the DB table name.
func (Extra) TableName() string {
    return "extras"
}

// NewExtra constructs a new Extra.
func NewExtra(userID uuid.UUID, category, text string) *Extra {
    return &Extra{
        ID:        uuid.New(),
        UserID:    userID,
        Category:  category,
        Text:      text,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }
}
