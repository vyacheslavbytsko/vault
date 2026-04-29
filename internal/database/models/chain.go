package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var chainNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

var (
	ErrInvalidChainName   = errors.New("chain name must contain only English letters, numbers, '_' or '-' characters")
	ErrChainOwnerRequired = errors.New("chain owner is required")
)

type Chain struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string     `gorm:"size:128;not null;uniqueIndex:idx_owner_name" json:"name"`
	OwnerID     uuid.UUID  `gorm:"column:owner_id;type:char(36);not null;index;uniqueIndex:idx_owner_name" json:"owner_id"`
	LastEventID *uuid.UUID `gorm:"column:last_event_id;type:char(36);index" json:"last_event_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`

	LastEvent *Event `gorm:"foreignKey:LastEventID;references:EventID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Owner     User   `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (c *Chain) TableName() string {
	return "chains"
}

func (c *Chain) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	return c.Validate()
}

func (c *Chain) Validate() error {
	if !chainNamePattern.MatchString(c.Name) {
		return ErrInvalidChainName
	}

	if c.OwnerID == uuid.Nil {
		return ErrChainOwnerRequired
	}

	return nil
}
