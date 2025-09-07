package domain_test

import (
	"testing"

	"github.com/brunojet/store-go/shared/internal/domain"
)

func TestContato_BeforeCreate_AllValid(t *testing.T) {
	c := domain.ApplicationContact{
		BaseEntity: domain.BaseEntity{
			Name: "Empresa X",
		},
		Site:  "www.empresax.com",
		Email: "contato@empresax.com",
		Phone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err != nil {
		t.Errorf("Esperado nil, recebeu erro: %v", err)
	}
}

func TestContato_BeforeCreate_RazaoSocialInvalida(t *testing.T) {
	c := domain.ApplicationContact{
		BaseEntity: domain.BaseEntity{
			Name: "X",
		},
		Site:  "www.empresax.com",
		Email: "contato@empresax.com",
		Phone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "razao_social deve ter pelo menos 3 caracteres" {
		t.Errorf("Esperado erro de razao_social, recebeu: %v", err)
	}
}

func TestContato_BeforeCreate_SiteInvalido(t *testing.T) {
	c := domain.ApplicationContact{
		BaseEntity: domain.BaseEntity{
			Name: "Empresa X",
		},
		Site:  "x",
		Email: "contato@empresax.com",
		Phone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "site deve ter pelo menos 3 caracteres" {
		t.Errorf("Esperado erro de site, recebeu: %v", err)
	}
}

func TestContato_BeforeCreate_EmailInvalido(t *testing.T) {
	c := domain.ApplicationContact{
		BaseEntity: domain.BaseEntity{
			Name: "Empresa X",
		},
		Site:  "www.empresax.com",
		Email: "x",
		Phone: "11999999999",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "email deve ter pelo menos 3 caracteres" {
		t.Errorf("Esperado erro de email, recebeu: %v", err)
	}
}

func TestContato_BeforeCreate_PhoneInvalido(t *testing.T) {
	c := domain.ApplicationContact{
		BaseEntity: domain.BaseEntity{
			Name: "Empresa X",
		},
		Site:  "www.empresax.com",
		Email: "contato@empresax.com",
		Phone: "1234567",
	}
	err := c.BeforeCreate(nil)
	if err == nil || err.Error() != "Phone deve ter pelo menos 8 caracteres" {
		t.Errorf("Esperado erro de Phone, recebeu: %v", err)
	}
}
