package domain

import "testing"

func TestImagem_TableName(t *testing.T) {
	var v Imagem
	want := "imagem"
	if got := v.TableName(); got != want {
		t.Fatalf("Imagem.TableName() = %q, want %q", got, want)
	}
}
