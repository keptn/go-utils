package models

import (
	"testing"
)

func TestIntegrationID_Hash(t *testing.T) {
	type fields struct {
		Name      string
		Namespace string
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
		},
			"34722d8839d1a4d9498467b96b1a4893a79bfd4b",
			false},
		{"missing name", fields{
			Name:      "",
			Namespace: "ns",
		},
			"",
			true},
		{"missing namespace", fields{
			Name:      "int",
			Namespace: "",
		},
			"",
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IntegrationID{
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
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
