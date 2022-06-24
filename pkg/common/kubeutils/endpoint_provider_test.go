package kubeutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestKeptnEndpointProvider_GetKeptnEndpointFromIngress_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "ingresses", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("error retrieving kubernetes ingress")
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromIngress(context.TODO(), "keptn", "ingress")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("error retrieving kubernetes ingress"), err)
}

func TestKeptnEndpointProvider_GetKeptnEndpointFromIngress_Invalid(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "ingresses", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1beta1.Ingress{Spec: v1beta1.IngressSpec{}}, nil
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromIngress(context.TODO(), "keptn", "ingress")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("cannot retrieve ingress data: ingress rule does not exist"), err)
}

func TestKeptnEndpointProvider_GetKeptnEndpointFromIngress_Valid(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "ingresses", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1beta1.Ingress{
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "1.1.1.1",
					},
				},
			},
		}, nil
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromIngress(context.TODO(), "keptn", "api-keptn-ingress")
	require.Equal(t, "1.1.1.1", res)
	require.Nil(t, err)
}

func TestKeptnEndpointProvider_GetKeptnEndpointFromService_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "services", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("error retrieving kubernetes service")
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromService(context.TODO(), "keptn", "service1")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("error retrieving kubernetes service"), err)
}

func TestKeptnEndpointProvider_GetKeptnEndpointFromService_NoLoadBalancer(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "services", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Service{Spec: v1.ServiceSpec{Type: v1.ServiceTypeClusterIP}}, nil
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromService(context.TODO(), "keptn22", "service")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("it doesn't support ClusterIP & NodePort type service for fetching endpoint automatically"), err)
}

func TestKeptnEndpointProvider_GetKeptnEndpointFromService_LoadBalancerNoIngress(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "services", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Service{
			Spec: v1.ServiceSpec{
				Type: v1.ServiceTypeLoadBalancer},
			Status: v1.ServiceStatus{
				LoadBalancer: v1.LoadBalancerStatus{},
			},
		}, nil
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromService(context.TODO(), "keptn22", "service")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("Loadbalancer IP isn't found"), err)
}

func TestKeptnEndpointProvider_GetKeptnEndpointFromService_LoadBalancerIngressIP(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "services", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Service{
			Spec: v1.ServiceSpec{Type: v1.ServiceTypeLoadBalancer},
			Status: v1.ServiceStatus{
				LoadBalancer: v1.LoadBalancerStatus{
					Ingress: []v1.LoadBalancerIngress{
						{IP: "1.1.1.1"},
					},
				},
			},
		}, nil
	})
	keptnEndpointProvider := &KeptnEndpointProvider{clientSet: kubernetes}
	res, err := keptnEndpointProvider.GetKeptnEndpointFromService(context.TODO(), "keptn22", "service")
	require.Equal(t, "1.1.1.1", res)
	require.Nil(t, err)
}
