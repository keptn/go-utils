package kubeutils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KeptnEndpointProvider wraps around the kubernetes interface to enhance testability
type KeptnEndpointProvider struct {
	clientSet kubernetes.Interface
}

// NewKeptnEndpointProvider creates new KeptnEndpointProvider
func NewKeptnEndpointProvider(useInClusterConfig bool) (*KeptnEndpointProvider, error) {
	clientSet, err := GetClientSet(useInClusterConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not create KeptnEndpointProvider: %s", err.Error())
	}
	return &KeptnEndpointProvider{clientSet: clientSet}, nil
}

// GetKeptnEndpointFromIngress returns the host of ingress object Keptn Installation
func (a *KeptnEndpointProvider) GetKeptnEndpointFromIngress(namespace string, ingressName string) (string, error) {
	keptnIngress, err := a.clientSet.ExtensionsV1beta1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if keptnIngress.Spec.Rules != nil {
		return keptnIngress.Spec.Rules[0].Host, nil
	}
	return "", fmt.Errorf("cannot retrieve ingress data: ingress rule does not exist")
}

// GetKeptnEndpointFromService returns the loadbalancer service IP from Keptn Installation
func (a *KeptnEndpointProvider) GetKeptnEndpointFromService(namespace string, serviceName string) (string, error) {
	keptnService, err := a.clientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	switch keptnService.Spec.Type {
	case "LoadBalancer":
		if len(keptnService.Status.LoadBalancer.Ingress) > 0 {
			return keptnService.Status.LoadBalancer.Ingress[0].IP, nil
		}
		return "", fmt.Errorf("Loadbalancer IP isn't found")
	default:
		return "", fmt.Errorf("It doesn't support ClusterIP & NodePort type service for fetching endpoint automatically")
	}
}
