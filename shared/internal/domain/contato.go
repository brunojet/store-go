package domain

import (
	"errors"

	"gorm.io/gorm"
)

type Contato struct {
	BaseEntity
	RazaoSocial string `gorm:"column:razao_social;not null"`
	Site        string `gorm:"column:site;not null"`
	Email       string `gorm:"column:email;not null"`
	Telefone    string `gorm:"column:telefone;not null"`
}

func (c *Contato) BeforeCreate(tx *gorm.DB) (err error) {
	if len(c.RazaoSocial) < 3 {
		return errors.New("razao_social deve ter pelo menos 3 caracteres")
	}
	if len(c.Site) < 3 {
		return errors.New("site deve ter pelo menos 3 caracteres")
	}
	if len(c.Email) < 3 {
		return errors.New("email deve ter pelo menos 3 caracteres")
	}
	if len(c.Telefone) < 8 {
		return errors.New("telefone deve ter pelo menos 8 caracteres")
	}
	return nil
}

func (Contato) TableName() string { return "cnt_app" }
