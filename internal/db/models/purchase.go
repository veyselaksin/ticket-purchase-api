package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Purchase struct {
	Id       string `gorm:"primaryKey"`
	TicketId string `gorm:"not null"`
	UserId   string `gorm:"not null"`
	Quantity int    `gorm:"not null"`

	// Relationships
	Ticket Ticket `gorm:"foreignKey:TicketId;references:Id"`

	// Audit fields
	CreatedBy string    `gorm:"not null"`
	UpdatedBy string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	IsActive  bool      `gorm:"default:true"`
}

// TableName specifies the table name for the Purchase model
func (Purchase) TableName() string {
	return "public.purchases"
}

// BeforeCreate is a GORM hook that is triggered before creating a new record
func (p *Purchase) BeforeCreate(tx *gorm.DB) error {
	p.Id = uuid.New().String()
	return nil
}
