package domain

type TerminalModel struct {
	BaseEntity
}

func (TerminalModel) TableName() string { return "terminal_model" }
