package kubeutils

import (
	"context"
	"fmt"

	typesv1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetKeptnEndpointFromIngress returns the host of ingress object Keptn Installation
func GetKeptnEndpointFromIngress(useInClusterConfig bool, namespace string, ingressName string) (string, error) {
	var keptnIngress *v1beta1.Ingress
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return "", err
	}
	keptnIngress, err = clientset.ExtensionsV1beta1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return keptnIngress.Spec.Rules[0].Host, nil
}

// GetKeptnEndpointFromService returns the loadbalancer service IP from Keptn Installation
func GetKeptnEndpointFromService(useInClusterConfig bool, namespace string, serviceName string) (string, error) {
	var keptnService *typesv1.Service
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return "", err
	}
	keptnService, err = clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
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
