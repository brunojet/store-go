package pkg

import (
	"fmt"
	"store-go/shared/internal"
)

func InitDBShared() {
	db, err := internal.InitDB()
	if err != nil {
		fmt.Println("Erro ao inicializar o banco:", err)
		return
	}
	fmt.Println("Banco inicializado com sucesso:", db)
}
