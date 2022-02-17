package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// KeptnContextExtendedCE keptn context extended CloudEvent
type KeptnContextExtendedCE struct {

	// contenttype
	Contenttype string `json:"contenttype,omitempty"`

	// data
	// Required: true
	Data interface{} `json:"data"`

	// extensions
	Extensions interface{} `json:"extensions,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// shkeptncontext
	Shkeptncontext string `json:"shkeptncontext,omitempty"`

	// shkeptnspecversion
	Shkeptnspecversion string `json:"shkeptnspecversion,omitempty"`

	// source
	// Required: true
	Source *string `json:"source"`

	// specversion
	Specversion string `json:"specversion,omitempty"`

	// time
	// Format: date-time
	Time time.Time `json:"time,omitempty"`

	// triggeredid
	Triggeredid string `json:"triggeredid,omitempty"`

	// gitcommitid
	GitCommitID string `json:"gitcommitid,omitempty"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// DataAs attempts to populate the provided data object with the event payload.
// data should be a pointer type.
func (ce *KeptnContextExtendedCE) DataAs(out interface{}) error {
	bytes, err := json.Marshal(ce.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

// Validate checks whether the required properties 'time', 'type', 'id' and 'source' are defined and non-empty
func (ce *KeptnContextExtendedCE) Validate() error {
	if ce.Time.IsZero() {
		return errors.New("time must be specified")
	}
	if ce.Type == nil || *ce.Type == "" {
		return errors.New("type must be specified")
	}
	if ce.ID == "" {
		return errors.New("id must be specified")
	}
	if ce.Source == nil || *ce.Source == "" {
		return errors.New("source must be specified")
	}
	return nil
}

// TemporaryData represents additional (temporary) data to be added
// to the data section of a keptn event
type TemporaryData interface{}

// AddTemporaryDataOptions are used to modify the behavior of adding temporary
// data to a keptn event
type AddTemporaryDataOptions struct {
	// OverwriteIfExisting indicates, that the data will be overwritten
	// in case a key for that data already exists
	OverwriteIfExisting bool
}

const temporaryDataRootKey = "temporaryData"

// AddTemporaryData adds further (temporary) properties to the data section of the keptn event
func (ce *KeptnContextExtendedCE) AddTemporaryData(key string, tmpData TemporaryData, opts AddTemporaryDataOptions) error {
	eventData := map[string]interface{}{}
	err := ce.DataAs(&eventData)
	if err != nil {
		return err
	}
	if temporaryData, found := eventData[temporaryDataRootKey]; found {
		if _, kfound := temporaryData.(map[string]interface{})[key]; kfound {
			if !opts.OverwriteIfExisting {
				return fmt.Errorf("Key %s already exists", key)
			}
			temporaryData.(map[string]interface{})[key] = tmpData
		}
	} else {
		eventData[temporaryDataRootKey] = map[string]interface{}{key: tmpData}
	}
	ce.Data = eventData
	return nil
}

// GetTemporaryData returns the (temporary) data eventually stored in the event
func (ce *KeptnContextExtendedCE) GetTemporaryData(key string, tmpdata interface{}) error {
	eventData := map[string]interface{}{}
	if err := ce.DataAs(&eventData); err != nil {
		return err
	}
	if temporaryData, found := eventData[temporaryDataRootKey]; found {
		if keyData, kfound := temporaryData.(map[string]interface{})[key]; kfound {
			if marshalledKeyData, err := json.Marshal(keyData); err == nil {
				return json.Unmarshal(marshalledKeyData, tmpdata)
			}
		}
		return fmt.Errorf("temporary data with key %s not found", key)
	}
	return fmt.Errorf("temporary data with key %s not found", key)
}

// ToJSON converts object to JSON string
func (ce *KeptnContextExtendedCE) ToJSON() ([]byte, error) {
	return json.Marshal(ce)
}

// FromJSON converts JSON string to object
func (ce *KeptnContextExtendedCE) FromJSON(b []byte) error {
	var res KeptnContextExtendedCE
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*ce = res
	return nil
}
