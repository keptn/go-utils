package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

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

// CheckDeploymentRolloutStatus checks the rollout status of the provided deployment
func CheckDeploymentRolloutStatus(serviceName string, namespace string) error {
	_, err := ExecuteCommand("kubectl", []string{"rollout", "status", "deployment/" + serviceName, "--namespace", namespace})
	return err
}

// GetKubeAPI returns the Kube API in version v1
func GetKubeAPI(useInClusterConfig bool) (v1.CoreV1Interface, error) {

	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
}
