package kubeutils

import (
	"context"
	"fmt"

	typesv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetKeptnAPITokenFromSecret returns the `keptn-api-token` data secret from Keptn Installation
func GetKeptnAPITokenFromSecret(useInClusterConfig bool, namespace string, secretName string) (string, error) {
	var keptnSecret *typesv1.Secret
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return "", err
	}
	keptnSecret, err = clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if apitoken, ok := keptnSecret.Data["keptn-api-token"]; ok {
		return string(apitoken), nil
	}
	return "", fmt.Errorf("data 'keptn-api-token' not found")
}
