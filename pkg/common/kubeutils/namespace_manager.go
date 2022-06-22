package kubeutils

import (
	"context"
	"fmt"

	typesv1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceManager struct {
	clientSet kubernetes.Interface
}

func NewManespaceManager(useInClusterConfig bool) (*NamespaceManager, error) {
	clientSet, err := GetClientSet(useInClusterConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not create ApiTokenProvider: %s", err.Error())
	}
	return &NamespaceManager{clientSet: clientSet}, nil
}

// ExistsNamespace checks whether a namespace with the provided name exists
func (a *NamespaceManager) ExistsNamespace(namespace string) (bool, error) {
	_, err := a.clientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*apierr.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateNamespace creates a new Kubernetes namespace with the provided name
func (a *NamespaceManager) CreateNamespace(namespace string, namespaceMetadata ...metav1.ObjectMeta) error {
	var buildNamespaceMetadata metav1.ObjectMeta
	if len(namespaceMetadata) > 0 {
		buildNamespaceMetadata = namespaceMetadata[0]
	}

	buildNamespaceMetadata.Name = namespace

	ns := &typesv1.Namespace{ObjectMeta: buildNamespaceMetadata}
	_, err := a.clientSet.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

// GetKeptnManagedNamespace returns the list of namespace with the annotation & label `keptn.sh/managed-by: keptn`
func (a *NamespaceManager) GetKeptnManagedNamespace() ([]string, error) {
	var namespaces []string

	namespaceList, err := a.clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "keptn.sh/managed-by",
	})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		if metav1.HasAnnotation(namespace.ObjectMeta, "keptn.sh/managed-by") {
			namespaces = append(namespaces, namespace.GetObjectMeta().GetName())
		}
	}
	return namespaces, nil
}
