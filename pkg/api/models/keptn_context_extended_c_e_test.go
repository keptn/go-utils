package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Gitcommitid        string
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
			ce := &models.KeptnContextExtendedCE{
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
				Gitcommitid:        tt.fields.Gitcommitid,
				Type:               tt.fields.Type,
			}
			if err := ce.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddTemporaryData(t *testing.T) {
	type TestData struct {
		v0_2_0.EventData
		Content string `json:"content"`
	}

	testData := TestData{
		EventData: v0_2_0.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
	}
	event, err := v0_2_0.KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
	require.Nil(t, err)

	type AdditionalData struct {
		SomeString string `json:"someString"`
		SomeInt    int    `json:"someInt"`
	}
	temporaryDataToAdd := models.TemporaryData(AdditionalData{
		SomeString: "Bernd",
		SomeInt:    34,
	})
	err = event.AddTemporaryData("distributor", temporaryDataToAdd, models.AddTemporaryDataOptions{})
	event.AddTemporaryData("distributor", temporaryDataToAdd, models.AddTemporaryDataOptions{})
	require.Nil(t, err)

	addi := AdditionalData{}
	err = event.GetTemporaryData("distributor", &addi)
	fmt.Println(addi.SomeInt)

	require.Nil(t, err)
}

func TestKeptnContextExtendedCE_TemporaryData(t *testing.T) {
	type TestData struct {
		v0_2_0.EventData
		Content string `json:"content"`
	}

	type AdditionalData struct {
		SomeString string `json:"someString"`
		SomeInt    int    `json:"someInt"`
	}

	testData := TestData{
		EventData: v0_2_0.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
	}
	t.Run("add temporary data without a key", func(t *testing.T) {
		ce, err := v0_2_0.KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		require.Nil(t, err)

		dataToAdd := AdditionalData{
			SomeString: "somestring",
			SomeInt:    2,
		}
		err = ce.AddTemporaryData("", dataToAdd, models.AddTemporaryDataOptions{})
		assert.Nil(t, err)
	})
	t.Run("add temporary data twice returns error", func(t *testing.T) {
		ce, err := v0_2_0.KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		require.Nil(t, err)

		dataToAdd := AdditionalData{
			SomeString: "somestring",
			SomeInt:    2,
		}
		err = ce.AddTemporaryData("the-key", dataToAdd, models.AddTemporaryDataOptions{})
		assert.Nil(t, err)

		dataRetrieved := AdditionalData{}
		err = ce.GetTemporaryData("the-key", &dataRetrieved)

		require.Nil(t, err)
		assert.Equal(t, dataToAdd, dataRetrieved)

		err = ce.AddTemporaryData("the-key", dataToAdd, models.AddTemporaryDataOptions{})
		assert.NotNil(t, err)
	})
	t.Run("add temporary data twice overwrites existing", func(t *testing.T) {
		ce, err := v0_2_0.KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		require.Nil(t, err)

		dataToAdd := AdditionalData{
			SomeString: "somestring",
			SomeInt:    2,
		}
		err = ce.AddTemporaryData("the-key", dataToAdd, models.AddTemporaryDataOptions{})
		assert.Nil(t, err)

		dataRetrieved := AdditionalData{}
		err = ce.GetTemporaryData("the-key", &dataRetrieved)

		require.Nil(t, err)
		assert.Equal(t, dataToAdd, dataRetrieved)

		dataToAdd.SomeInt = 1

		err = ce.AddTemporaryData("the-key", dataToAdd, models.AddTemporaryDataOptions{OverwriteIfExisting: true})
		assert.Nil(t, err)
		dataRetrieved = AdditionalData{}
		err = ce.GetTemporaryData("the-key", &dataRetrieved)
		require.Nil(t, err)
		assert.Equal(t, dataToAdd, dataRetrieved)
	})
	t.Run("get non existing temporary data", func(t *testing.T) {
		ce, err := v0_2_0.KeptnEvent("sh.keptn.event.dev.delivery.triggered", "source", testData).Build()
		require.Nil(t, err)

		dataRetrieved := AdditionalData{}
		err = ce.GetTemporaryData("the-key", &dataRetrieved)
		assert.NotNil(t, err)
	})
}
