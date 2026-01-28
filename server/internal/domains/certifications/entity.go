package certifications

import (
    "time"

    "github.com/google/uuid"
)

// Certification represents a user's certification or credential.
type Certification struct {
    ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
    UserID      uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
    Name        string    `gorm:"column:name;size:255;not null" json:"name"`
    Issuer      string    `gorm:"column:issuer;size:255" json:"issuer,omitempty"`
    Date        string    `gorm:"column:date;size:32" json:"date,omitempty"`
    URL         string    `gorm:"column:url;size:512" json:"url,omitempty"`
    Description string    `gorm:"column:description;type:text" json:"description,omitempty"`
    CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName sets the DB table name.
func (Certification) TableName() string {
    return "certifications"
}

// NewCertification constructs a new Certification.
func NewCertification(userID uuid.UUID, name string) *Certification {
    return &Certification{
        ID:        uuid.New(),
        UserID:    userID,
        Name:      name,
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
    }
}
