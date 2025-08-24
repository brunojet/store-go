package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type BaseEntity struct {
	BaseModel
	Nome      string `gorm:"not null"` // Obrigatório
	Descricao string
	Ativo     bool
}

func (e *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if len(e.Nome) < 1 {
		return errors.New("nome deve ter pelo menos 1 caracter")
	}
	return nil
}

type Anexo struct {
	BaseModel
	Nome          string `gorm:"size:128;not null"`
	TipoMime      string `gorm:"size:64;not null"`
	MD5           string `gorm:"size:32;not null"`
	Tamanho       int64  `gorm:"not null"`
	Armazenamento string `gorm:"size:256;not null"`
	Caminho       string `gorm:"size:256;not null"`
	Presente      bool   `gorm:"default:false"`
}

func (Anexo) TableName() string { return "anexo" }
