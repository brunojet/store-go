package domain_test

import (
	"testing"

	"github.com/brunojet/store-go/shared/internal/domain"
)

func TestContato_BeforeCreate_AllValid(t *testing.T) {
	c := domain.ContatoAplicativo{
		BaseEntity: domain.BaseEntity{
			Nome: "Empresa X",
		},
		Site:     "www.empresax.com",
		Email:    "contato@empresax.com",
		Telefone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err != nil {
		t.Errorf("Esperado nil, recebeu erro: %v", err)
	}
}

func TestContato_BeforeCreate_RazaoSocialInvalida(t *testing.T) {
	c := domain.ContatoAplicativo{
		BaseEntity: domain.BaseEntity{
			Nome: "X",
		},
		Site:     "www.empresax.com",
		Email:    "contato@empresax.com",
		Telefone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "razao_social deve ter pelo menos 3 caracteres" {
		t.Errorf("Esperado erro de razao_social, recebeu: %v", err)
	}
}

func TestContato_BeforeCreate_SiteInvalido(t *testing.T) {
	c := domain.ContatoAplicativo{
		BaseEntity: domain.BaseEntity{
			Nome: "Empresa X",
		},
		Site:     "x",
		Email:    "contato@empresax.com",
		Telefone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "site deve ter pelo menos 3 caracteres" {
		t.Errorf("Esperado erro de site, recebeu: %v", err)
	}
}

func TestContato_BeforeCreate_EmailInvalido(t *testing.T) {
	c := domain.ContatoAplicativo{
		BaseEntity: domain.BaseEntity{
			Nome: "Empresa X",
		},
		Site:     "www.empresax.com",
		Email:    "x",
		Telefone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "email deve ter pelo menos 3 caracteres" {
		t.Errorf("Esperado erro de email, recebeu: %v", err)
	}
}

func TestContato_BeforeCreate_TelefoneInvalido(t *testing.T) {
	c := domain.ContatoAplicativo{
		BaseEntity: domain.BaseEntity{
			Nome: "Empresa X",
		},
		Site:     "www.empresax.com",
		Email:    "contato@empresax.com",
		Telefone: "1234567",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "telefone deve ter pelo menos 8 caracteres" {
		t.Errorf("Esperado erro de telefone, recebeu: %v", err)
	}
}
