package kubeutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestAPITokenProvider_GetKeptnAPITokenFromSecret_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("error retrieving kubernetes secret")
	})
	apiTokenProvider := &APITokenProvider{clientSet: kubernetes}
	res, err := apiTokenProvider.GetKeptnAPITokenFromSecret(context.TODO(), "keptn", "secret")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("error retrieving kubernetes secret"), err)

}

func TestAPITokenProvider_GetKeptnAPITokenFromSecret_InvalidData(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{Data: map[string][]byte{"some-data": []byte("token")}}, nil
	})
	apiTokenProvider := &APITokenProvider{clientSet: kubernetes}
	res, err := apiTokenProvider.GetKeptnAPITokenFromSecret(context.TODO(), "keptn", "secret")
	require.Equal(t, "", res)
	require.Error(t, err)
	require.Equal(t, fmt.Errorf("data 'keptn-api-token' not found"), err)

}

func TestAPITokenProvider_GetKeptnAPITokenFromSecret_ValidData(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{Data: map[string][]byte{"keptn-api-token": []byte("token")}}, nil
	})
	apiTokenProvider := &APITokenProvider{clientSet: kubernetes}
	res, err := apiTokenProvider.GetKeptnAPITokenFromSecret(context.TODO(), "keptn", "secret")
	require.Equal(t, "token", res)
	require.Nil(t, err)

}
