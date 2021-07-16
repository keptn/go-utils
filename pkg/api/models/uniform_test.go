package models

import (
	"testing"
)

func TestIntegrationID_Hash(t *testing.T) {
	type fields struct {
		Name      string
		Namespace string
		NodeName  string
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
			NodeName:  "n1",
		},
			"ea4caa30233315068a5d2e3e7bb480851379c1e9",
			false},
		{"missing name", fields{
			Name:      "",
			Namespace: "ns",
			NodeName:  "n1",
		},
			"",
			true},
		{"missing namespace", fields{
			Name:      "int",
			Namespace: "",
			NodeName:  "n1",
		},
			"",
			true},
		{"missing nodename", fields{
			Name:      "int",
			Namespace: "ns",
			NodeName:  "",
		},
			"",
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IntegrationID{
				Name:      tt.fields.Name,
				Namespace: tt.fields.Namespace,
				NodeName:  tt.fields.NodeName,
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
