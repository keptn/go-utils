package models

import (
	"testing"
)

func TestIntegrationID_Hash(t *testing.T) {
	type fields struct {
		Name      string
		Namespace string
		Project   string
		Stage     string
		Service   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"ok", fields{
			Name:      "int",
			Namespace: "ns",
			Project:   "pr",
			Stage:     "st",
			Service:   "svc",
		},
			"f8f6ee30ef9330dc4c6c6167cd1434ca17d54c73",
			false},
		{"missing name", fields{
			Name:      "",
			Namespace: "ns",
			Project:   "pr",
			Stage:     "st",
			Service:   "svc",
		},
			"",
			true},
		{"missing namespace", fields{
			Name:      "int",
			Namespace: "",
			Project:   "pr",
			Stage:     "st",
			Service:   "svc",
		},
			"",
			true},
		{"missing project", fields{
			Name:      "int",
			Namespace: "ns",
			Project:   "",
			Stage:     "st",
			Service:   "svc",
		},
			"49a6b635c8a4588db5986bf888126bc56238c9df",
			false},
		{"missing stage", fields{
			Name:      "int",
			Namespace: "ns",
			Project:   "pr",
			Stage:     "",
			Service:   "svc",
		},
			"d7065a3f06078d28f00e44733d840a32d0ee2a07",
			false},
		{"missing service", fields{
			Name:      "int",
			Namespace: "ns",
			Project:   "pr",
			Stage:     "st",
			Service:   "",
		},
			"3ce10a2ca2c2d5a69f42cf9cd1e392186c124c8e",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IntegrationID{
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
				Project:   tt.fields.Project,
				Stage:     tt.fields.Stage,
				Service:   tt.fields.Service,
			}
			got, err := i.Hash()
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
