package domain

type TipoIntegracao struct {
	BaseEntity
}

func (TipoIntegracao) TableName() string { return "tip_itgr" }
