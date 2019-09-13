package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"

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
	_, err = ExecuteCommand("helm", []string{"upgrade", "--install", projectStage, helmChart, "--namespace", projectStage, "--wait"})
	return err
}

// RestartPodsWithSelector restarts the pods which are found in the provided namespace and selector
func RestartPodsWithSelector(useInClusterConfig bool, namespace string, selector string) error {
	clientset, err := GetKubeAPI(useInClusterConfig)
	if err != nil {
		return err
	}
	pods, err := clientset.Pods(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if err := clientset.Pods(namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func ScaleDeployment(useInClusterConfig bool, deployment string, namespace string, replicas int32) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(deployment, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}

		result.Spec.Replicas = int32Ptr(replicas)
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	return retryErr
}

func int32Ptr(i int32) *int32 { return &i }

// WaitForDeploymentToBeRolledOut waits until the deployment is Available
func WaitForDeploymentToBeRolledOut(useInClusterConfig bool, deploymentName string, namespace string) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}

	deployment, err := getDeployment(clientset, namespace, deploymentName)
	for deployment.Generation <= deployment.Status.ObservedGeneration {

		var cond *appsv1.DeploymentCondition

		for i := range deployment.Status.Conditions {
			c := deployment.Status.Conditions[i]
			if c.Type == appsv1.DeploymentProgressing {
				cond = &c
				break
			}
		}

		if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
			return fmt.Errorf("Deployment %q exceeded its progress deadline", deployment.Name)
		}
		if !(deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas ||
			deployment.Status.Replicas > deployment.Status.UpdatedReplicas ||
			deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas) {
			return nil
		}

		time.Sleep(2 * time.Second)
		deployment, err = getDeployment(clientset, namespace, deploymentName)
		if err != nil {
			return err
		}
	}

	return nil
}

// WaitForDeploymentsInNamespace waits until all deployments in a namespace are available
func WaitForDeploymentsInNamespace(useInClusterConfig bool, namespace string) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}
	deps, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	for _, dep := range deps.Items {
		WaitForDeploymentToBeRolledOut(useInClusterConfig, dep.Name, namespace)
	}
	return nil
}

func getDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	dep, err := clientset.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil &&
		strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
		time.Sleep(10 * time.Second)
		return clientset.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
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

// GetKeptnDomain reads the configmap keptn-domain in namespace keptn and returns
// the contained app_domain
func GetKeptnDomain(useInClusterConfig bool) (string, error) {
	api, err := GetKubeAPI(useInClusterConfig)
	if err != nil {
		return "", err
	}

	cm, err := api.ConfigMaps("keptn").Get("keptn-domain", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return cm.Data["app_domain"], nil
}
