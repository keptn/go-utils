package models

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Integration struct {
	ID       string   `json:"id" bson:"_id"`
	Name     string   `json:"name" bson:"name"`
	MetaData MetaData `json:"metadata" bson:"metadata"`
	// Deprecated: for backwards compatibility Subscription is populated
	// but new code shall use Subscriptions
	Subscription  Subscription        `json:"subscription" bson:"subscription"`
	Subscriptions []TopicSubscription `json:"subscriptions" bson:"subscriptions"`
}

type MetaData struct {
	Hostname           string             `json:"hostname" bson:"hostname"`
	IntegrationVersion string             `json:"integrationversion" bson:"integrationversion"`
	DistributorVersion string             `json:"distributorversion" bson:"distributorversion"`
	Location           string             `json:"location" bson:"location"`
	KubernetesMetaData KubernetesMetaData `json:"kubernetesmetadata" bson:"kubernetesmetadata"`
	LastSeen           time.Time          `json:"lastseen" bson:"lastseen"`
}

// Deprecated
type Subscription struct {
	Topics []string           `json:"topics" bson:"topics"`
	Status string             `json:"status" bson:"status"`
	Filter SubscriptionFilter `json:"filter" bson:"filter"`
}

type TopicSubscription struct {
	Topics []string                `json:"topics" bson:"topics"`
	Status string                  `json:"status" bson:"status"`
	Filter TopicSubscriptionFilter `json:"filter" bson:"filter"`
}

// Deprecated
type SubscriptionFilter struct {
	// Deprecated: for backwards compatibility Project is still populated
	// but new code shall use Projects
	Project string `json:"project" bson:"project"`
	// Deprecated: for backwards compatibility Stage is still populated
	// but new code shall use Stages
	Stage string `json:"stage" bson:"stage"`
	// Deprecated: for backwards compatibility Service is still populated
	// but new code shall use Services
	Service string `json:"service" bson:"service"`
}

type TopicSubscriptionFilter struct {
	Projects []string `json:"projects" bson:"projects"`
	Stages   []string `json:"stages" bson:"stages"`
	Services []string `json:"services" bson:"services"`
}

type KubernetesMetaData struct {
	Namespace      string `json:"namespace" bson:"namespace"`
	PodName        string `json:"podname" bson:"podname"`
	DeploymentName string `json:"deploymentname" bson:"deploymentname"`
}

type IntegrationID struct {
	Name      string `json:"name" bson:"name"`
	Namespace string `json:"namespace" bson:"namespace"`
	NodeName  string `json:"nodename" bson:"nodename"`
}

func (i IntegrationID) Hash() (string, error) {
	if !i.validate() {
		return "", fmt.Errorf("incomplete integration ID. At least 'name','namespace' and 'nodename' must be set")
	}
	raw := fmt.Sprintf("%s-%s-%s", i.Name, i.Namespace, i.NodeName)
	hasher := sha1.New() //nolint:gosec
	hasher.Write([]byte(raw))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (i IntegrationID) validate() bool {
	return i.Name != "" && i.Namespace != "" && i.NodeName != ""
}

func (i *Integration) ToJSON() ([]byte, error) {
	if i == nil {
		return nil, nil
	}
	return json.Marshal(i)
}

func (i *Integration) FromJSON(b []byte) error {
	var res Integration
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*i = res
	return nil
}
