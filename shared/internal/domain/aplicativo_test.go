package domain

import "testing"

func TestAplicativo_TableName_And_BeforeCreate(t *testing.T) {
	var a Aplicativo
	if got := a.TableName(); got == "" {
		t.Fatalf("TableName() returned empty string")
	}

	cases := []struct {
		name    string
		fill    func(a *Aplicativo)
		wantErr bool
	}{
		{"empty", func(a *Aplicativo) {}, true},
		{"only razao short", func(a *Aplicativo) { a.RazaoSocial = "AB" }, true},
		{"razao ok site short", func(a *Aplicativo) { a.RazaoSocial = "ABC"; a.Site = "ab" }, true},
		{"razao site ok email short", func(a *Aplicativo) { a.RazaoSocial = "ABC"; a.Site = "ABC"; a.Email = "e" }, true},
		{"telefone short", func(a *Aplicativo) { a.RazaoSocial = "ABC"; a.Site = "ABC"; a.Email = "e@e.com"; a.Telefone = "123" }, true},
		{"all ok", func(a *Aplicativo) {
			a.RazaoSocial = "ABC"
			a.Site = "ABC"
			a.Email = "e@e.com"
			a.Telefone = "12345678"
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var x Aplicativo
			tc.fill(&x)
			err := x.BeforeCreate(nil)
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Fatalf("case %s: gotErr=%v wantErr=%v (err=%v)", tc.name, gotErr, tc.wantErr, err)
			}
		})
	}
}
