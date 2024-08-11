package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Ticket struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Allocation  int    `json:"allocation" gorm:"not null"`

	// Audit fields
	CreatedBy string    `json:"created_by" gorm:"not null"`
	UpdatedBy string    `json:"updated_by" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
}

// TableName specifies the table name for the Ticket model
func (Ticket) TableName() string {
	return "public.tickets"
}

// BeforeCreate is a GORM hook that is triggered before creating a new record
func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	t.Id = uuid.New().String()
	return nil
}
