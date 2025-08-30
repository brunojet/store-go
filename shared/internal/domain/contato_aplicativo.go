package domain

import (
	"errors"

	"gorm.io/gorm"
)

type ContatoAplicativo struct {
	BaseEntity
	Site     string `gorm:"column:site;not null"`
	Email    string `gorm:"column:email;not null"`
	Telefone string `gorm:"column:telefone;not null"`
}

func (c *ContatoAplicativo) BeforeCreate(tx *gorm.DB) (err error) {
	if len(c.Nome) < 3 {
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

func (ContatoAplicativo) TableName() string { return "cnt_aplv" }
