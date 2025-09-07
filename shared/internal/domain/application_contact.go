package domain

import (
	"errors"

	"gorm.io/gorm"
)

type ApplicationContact struct {
	BaseEntity
	Site  string `gorm:"not null"`
	Email string `gorm:"not null"`
	Phone string `gorm:"not null"`
}

func (c *ApplicationContact) BeforeCreate(tx *gorm.DB) (err error) {
	if len(c.Name) < 3 {
		return errors.New("name deve ter pelo menos 3 caracteres")
	}
	if len(c.Site) < 3 {
		return errors.New("site deve ter pelo menos 3 caracteres")
	}
	if len(c.Email) < 3 {
		return errors.New("email deve ter pelo menos 3 caracteres")
	}
	if len(c.Phone) < 8 {
		return errors.New("phone deve ter pelo menos 8 caracteres")
	}
	return nil
}

func (ApplicationContact) TableName() string { return "application_contact" }
