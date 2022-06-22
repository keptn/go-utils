package kubeutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestApiTokenProvider_GetKeptnAPITokenFromSecret_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("Error retrieving kubernetes secret")
	})
	apiTokenProvider := &ApiTokenProvider{clientSet: kubernetes}
	res, err := apiTokenProvider.GetKeptnAPITokenFromSecret("keptn", "secret")
	require.Equal(t, "", res)
	require.Error(t, err)

}

func TestApiTokenProvider_GetKeptnAPITokenFromSecret_InvalidData(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{Data: map[string][]byte{"some-data": []byte("token")}}, nil
	})
	apiTokenProvider := &ApiTokenProvider{clientSet: kubernetes}
	res, err := apiTokenProvider.GetKeptnAPITokenFromSecret("keptn", "secret")
	require.Equal(t, "", res)
	require.Error(t, err)

}

func TestApiTokenProvider_GetKeptnAPITokenFromSecret_ValidData(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "secrets", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Secret{Data: map[string][]byte{"keptn-api-token": []byte("token")}}, nil
	})
	apiTokenProvider := &ApiTokenProvider{clientSet: kubernetes}
	res, err := apiTokenProvider.GetKeptnAPITokenFromSecret("keptn", "secret")
	require.Equal(t, "token", res)
	require.Nil(t, err)

}
