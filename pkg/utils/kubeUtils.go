package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// DoHelmUpgrade executes a helm update and upgrade
func DoHelmUpgrade(project string, stage string) error {
	helmChart := fmt.Sprintf("%s/helm-chart", project)
	projectStage := fmt.Sprintf("%s-%s", project, stage)
	_, err := ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		return err
	}
	_, err = ExecuteCommand("helm", []string{"dep", "update", helmChart})
	if err != nil {
		return err
	}
	_, err = ExecuteCommand("helm", []string{"upgrade", "--install", projectStage, helmChart, "--namespace", projectStage})
	return err
}

// WaitForDeploymentToBeAvailable waits until the deployment is Available
func WaitForDeploymentToBeAvailable(useInClusterConfig bool, serviceName string, namespace string) error {

	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}

	dep, err := getDeployment(clientset, namespace, serviceName)

	for dep.Status.UnavailableReplicas > 0 {
		time.Sleep(2 * time.Second)
		dep, err = getDeployment(clientset, namespace, serviceName)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDeployment(clientset *kubernetes.Clientset, namespace string, serviceName string) (*appsv1.Deployment, error) {
	dep, err := clientset.AppsV1().Deployments(namespace).Get(serviceName, metav1.GetOptions{})
	if err != nil &&
		strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
		time.Sleep(10 * time.Second)
		return clientset.AppsV1().Deployments(namespace).Get(serviceName, metav1.GetOptions{})
	}
	return dep, nil
}

// GetKubeAPI returns the CoreV1Interface
func GetKubeAPI(useInClusterConfig bool) (v1.CoreV1Interface, error) {

	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
}

// GetClientset returns the kubernetes Clientset
func GetClientset(useInClusterConfig bool) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(
			UserHomeDir(), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
