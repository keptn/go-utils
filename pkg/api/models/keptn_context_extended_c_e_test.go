package models

import (
	"testing"
	"time"
)

func TestKeptnContextExtendedCE_Validate(t *testing.T) {
	source := "my-source"
	eventType := "my-type"
	type fields struct {
		Contenttype        string
		Data               interface{}
		Extensions         interface{}
		ID                 string
		Shkeptncontext     string
		Shkeptnspecversion string
		Source             *string
		Specversion        string
		Time               time.Time
		Triggeredid        string
		Type               *string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validation passes",
			fields: fields{
				ID:     "my-id",
				Source: &source,
				Time:   time.Now(),
				Type:   &eventType,
			},
			wantErr: false,
		},
		{
			name: "missing type",
			fields: fields{
				ID:     "my-id",
				Source: &source,
				Time:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing id",
			fields: fields{
				Source: &source,
				Time:   time.Now(),
				Type:   &eventType,
			},
			wantErr: true,
		},
		{
			name: "missing time",
			fields: fields{
				ID:     "my-id",
				Source: &source,
				Type:   &eventType,
			},
			wantErr: true,
		},
		{
			name: "missing source",
			fields: fields{
				ID:   "my-id",
				Time: time.Now(),
				Type: &eventType,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := &KeptnContextExtendedCE{
				Contenttype:        tt.fields.Contenttype,
				Data:               tt.fields.Data,
				Extensions:         tt.fields.Extensions,
				ID:                 tt.fields.ID,
				Shkeptncontext:     tt.fields.Shkeptncontext,
				Shkeptnspecversion: tt.fields.Shkeptnspecversion,
				Source:             tt.fields.Source,
				Specversion:        tt.fields.Specversion,
				Time:               tt.fields.Time,
				Triggeredid:        tt.fields.Triggeredid,
				Type:               tt.fields.Type,
			}
			if err := ce.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
