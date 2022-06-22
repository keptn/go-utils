package kubeutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	typesv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestNamespaceManager_ExistsNamespace_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("Error retrieving kubernetes namespaces")
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	res, err := namespaceManager.ExistsNamespace("keptn")
	require.Equal(t, false, res)
	require.Error(t, err)
}

func TestNamespaceManager_ExistsNamespace_NotExists(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("some err")
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	res, err := namespaceManager.ExistsNamespace("keptn")
	require.Equal(t, false, res)
	require.Nil(t, err)
}

func TestNamespaceManager_ExistsNamespace_Exists(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("get", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &typesv1.Namespace{}, nil
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	res, err := namespaceManager.ExistsNamespace("keptn")
	require.Equal(t, true, res)
	require.Nil(t, err)
}

func TestNamespaceManager_CreateNamespace_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("create", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("Error creating kubernetes namespace")
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	err := namespaceManager.CreateNamespace("keptn")
	require.Error(t, err)
}

func TestNamespaceManager_CreateNamespace_Success(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("create", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &typesv1.Namespace{}, nil
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	err := namespaceManager.CreateNamespace("keptn")
	require.Error(t, err)
}

func TestNamespaceManager_GetKeptnManagedNamespace_FailClientSet(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("create", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("Error retrieving namespaces")
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	res, err := namespaceManager.GetKeptnManagedNamespace()
	require.Equal(t, []string([]string(nil)), res)
	require.Error(t, err)
}

func TestNamespaceManager_GetKeptnManagedNamespace_Success(t *testing.T) {
	kubernetes := fake.NewSimpleClientset()
	kubernetes.Fake.PrependReactor("create", "namespace", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &typesv1.NamespaceList{
			Items: []typesv1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"keptn.sh/managed-by": "keptn.sh/managed-by"}}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"keptn.sh/managed-by": "keptn.sh/managed-by"}}},
			},
		}, nil
	})
	namespaceManager := &NamespaceManager{clientSet: kubernetes}
	res, err := namespaceManager.GetKeptnManagedNamespace()
	require.Equal(t, []string([]string(nil)), res)
	require.Nil(t, err)
}
