package domain

type ModeloTerminal struct {
	BaseEntity
}

func (ModeloTerminal) TableName() string { return "mdl_trml" }
