package models

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Integration represents a Keptn service a.k.a. Keptn sntegration
// and contains the name, id and subscription data as well as other information
// needed to register a Keptn service to the control plane
type Integration struct {
	ID       string   `json:"id" bson:"_id"`
	Name     string   `json:"name" bson:"name"`
	MetaData MetaData `json:"metadata" bson:"metadata"`
	// Deprecated: for backwards compatibility Subscription is populated
	// but new code shall use Subscriptions
	Subscription  Subscription        `json:"subscription" bson:"subscription"`
	Subscriptions []EventSubscription `json:"subscriptions" bson:"subscriptions"`
}

// MetaData contains important information about the Keptn service which is used
// during registering the service to the control plane
type MetaData struct {
	Hostname           string             `json:"hostname" bson:"hostname"`
	IntegrationVersion string             `json:"integrationversion" bson:"integrationversion"`
	DistributorVersion string             `json:"distributorversion" bson:"distributorversion"`
	Location           string             `json:"location" bson:"location"`
	KubernetesMetaData KubernetesMetaData `json:"kubernetesmetadata" bson:"kubernetesmetadata"`
	LastSeen           time.Time          `json:"lastseen" bson:"lastseen"`
}

// Subscription describes to what events the Keptn service is subscribed to
// Deprecated
type Subscription struct {
	Topics []string           `json:"topics" bson:"topics"`
	Status string             `json:"status" bson:"status"`
	Filter SubscriptionFilter `json:"filter" bson:"filter"`
}

// EventSubscription describes to what events the Keptn service is subscribed to
type EventSubscription struct {
	ID     string                  `json:"id" bson:"id"`
	Event  string                  `json:"event" bson:"event"`
	Filter EventSubscriptionFilter `json:"filter" bson:"filter"`
}

// SubscriptionFilter is used to filter subscriptions by project stage or service
// Deprecated
type SubscriptionFilter struct {
	Project string `json:"project" bson:"project"`
	Stage   string `json:"stage" bson:"stage"`
	Service string `json:"service" bson:"service"`
}

// EventSubscriptionFilter is used to filter subscriptions by projects stages and/or services
type EventSubscriptionFilter struct {
	Projects []string `json:"projects" bson:"projects"`
	Stages   []string `json:"stages" bson:"stages"`
	Services []string `json:"services" bson:"services"`
}

// KubernetesMetaData represents metadata specific to Kubernetes
type KubernetesMetaData struct {
	Namespace      string `json:"namespace" bson:"namespace"`
	PodName        string `json:"podname" bson:"podname"`
	DeploymentName string `json:"deploymentname" bson:"deploymentname"`
}

// IntegrationID is the unique id of a Keptn service a.k.a "Keptn integration"
// It is composed by a name, the namespace the service resides in and the node name of the cluster node
type IntegrationID struct {
	Name      string `json:"name" bson:"name"`
	Namespace string `json:"namespace" bson:"namespace"`
	NodeName  string `json:"nodename" bson:"nodename"`
}

// Hash computes a hash value of an IntegrationID
// The IntegrationID must have a name, namespace as well as a nodename set
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

func (s *EventSubscription) ToJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}
