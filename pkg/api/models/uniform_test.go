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
			"",
			true},
		{"missing stage", fields{
			Name:      "int",
			Namespace: "ns",
			Project:   "pr",
			Stage:     "",
			Service:   "svc",
		},
			"",
			true},
		{"missing service", fields{
			Name:      "int",
			Namespace: "ns",
			Project:   "pr",
			Stage:     "st",
			Service:   "",
		},
			"",
			true},
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
