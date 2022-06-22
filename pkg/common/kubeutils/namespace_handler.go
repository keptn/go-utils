package kubeutils

import (
	"context"

	typesv1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ExistsNamespace checks whether a namespace with the provided name exists
func ExistsNamespace(useInClusterConfig bool, namespace string) (bool, error) {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return false, err
	}
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*apierr.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateNamespace creates a new Kubernetes namespace with the provided name
func CreateNamespace(useInClusterConfig bool, namespace string, namespaceMetadata ...metav1.ObjectMeta) error {

	var buildNamespaceMetadata metav1.ObjectMeta
	if len(namespaceMetadata) > 0 {
		buildNamespaceMetadata = namespaceMetadata[0]
	}

	buildNamespaceMetadata.Name = namespace

	ns := &typesv1.Namespace{ObjectMeta: buildNamespaceMetadata}
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

// PatchKeptnManagedNamespace to patch the namespace with the annotation & label `keptn.sh/managed-by: keptn`
func PatchKeptnManagedNamespace(useInClusterConfig bool, namespace string) error {
	var patchData = []byte(`{"metadata": {"annotations": {"keptn.sh/managed-by": "keptn"}, "labels": {"keptn.sh/managed-by": "keptn"}}}`)
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().Namespaces().Patch(context.TODO(), namespace, types.StrategicMergePatchType, patchData,
		metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetKeptnManagedNamespace returns the list of namespace with the annotation & label `keptn.sh/managed-by: keptn`
func GetKeptnManagedNamespace(useInClusterConfig bool) ([]string, error) {
	var namespaceList *typesv1.NamespaceList
	var namespaces []string
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}
	namespaceList, err = clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
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
