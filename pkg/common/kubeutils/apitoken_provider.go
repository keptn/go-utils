package kubeutils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ApiTokenProvider struct {
	clientSet kubernetes.Interface
}

// NewApiTokenProvider creates new ApiTokenProvider
func NewApiTokenProvider(useInClusterConfig bool) (*ApiTokenProvider, error) {
	clientSet, err := GetClientSet(useInClusterConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not create ApiTokenProvider: %s", err.Error())
	}
	return &ApiTokenProvider{clientSet: clientSet}, nil
}

// GetKeptnAPITokenFromSecret returns the `keptn-api-token` data secret from Keptn Installation
func (a *ApiTokenProvider) GetKeptnAPITokenFromSecret(namespace string, secretName string) (string, error) {
	keptnSecret, err := a.clientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if apitoken, ok := keptnSecret.Data["keptn-api-token"]; ok {
		return string(apitoken), nil
	}
	return "", fmt.Errorf("data 'keptn-api-token' not found")
}
