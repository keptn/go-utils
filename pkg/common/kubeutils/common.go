package kubeutils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"

	typesv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	apierr "k8s.io/apimachinery/pkg/api/errors"

	// Initialize all known client auth plugins.
	_ "github.com/Azure/go-autorest/autorest"

	// Initialize all known client auth plugins.
	appsv1 "k8s.io/api/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const keptnFolderName = ".keptn"

// UserHomeDir returns the HOME directory by taking into account the operating system
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// ExpandTilde expands ~ to HOME
func ExpandTilde(fileName string) string {
	if fileName == "~" {
		return UserHomeDir()
	} else if strings.HasPrefix(fileName, "~/") {
		return filepath.Join(UserHomeDir(), fileName[2:])
	}
	return fileName
}

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

// GetClientset returns the kubernetes Clientset
func GetClientset(useInClusterConfig bool) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		var kubeconfig string
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = ExpandTilde(os.Getenv("KUBECONFIG"))
		} else {
			kubeconfig = filepath.Join(
				UserHomeDir(), ".kube", "config",
			)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// GetKeptnAPITokenFromSecret returns the `keptn-api-token` data secret from Keptn Installation
func GetKeptnAPITokenFromSecret(useInClusterConfig bool, namespace string, secretName string) (string, error) {
	var keptnSecret *typesv1.Secret
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return "", err
	}
	keptnSecret, err = clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if apitoken, ok := keptnSecret.Data["keptn-api-token"]; ok {
		return string(apitoken), nil
	}
	return "", fmt.Errorf("data 'keptn-api-token' not found")
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

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

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

// GetKeptnDirectory returns a path, which is used to store logs and possibly creds
func GetKeptnDirectory() (string, error) {

	keptnDir := UserHomeDir() + string(os.PathSeparator) + keptnFolderName + string(os.PathSeparator)

	if _, err := os.Stat(keptnDir); os.IsNotExist(err) {
		err := os.MkdirAll(keptnDir, os.ModePerm)
		fmt.Println("keptn creates the folder " + keptnDir + " to store logs and possibly creds.")
		if err != nil {
			return "", err
		}
	}

	return keptnDir, nil
}

// IsDeployment tests whether the provided struct is a deployment
func IsDeployment(dpl *appsv1.Deployment) bool {
	return strings.ToLower(dpl.Kind) == "deployment"
}

// WaitForDeploymentToBeRolledOut waits until the deployment is Available
func WaitForDeploymentToBeRolledOut(useInClusterConfig bool, deploymentName string, namespace string) error {
	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return err
	}

	const maxWaitForDeploymentRetries = 90
	deployment, err := getDeployment(clientset, namespace, deploymentName)
	retries := 0
	for {

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
		retries = retries + 1
		if retries >= maxWaitForDeploymentRetries {
			return fmt.Errorf("Timed out waiting for deployment %q", deployment.Name)
		}
	}
}

func getDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	dep, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil &&
		strings.Contains(err.Error(), "the object has been modified; please apply your changes to the latest version and try again") {
		time.Sleep(10 * time.Second)
		return clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	}
	return dep, nil
}

// IsService tests whether the provided struct is a service
func IsService(svc *typesv1.Service) bool {
	return strings.ToLower(svc.Kind) == "service"
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

// ExecuteCommandInDirectory executes the command using the args within the specified directory
func ExecuteCommandInDirectory(command string, args []string, directory string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = directory
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

func GetHelmChartURI(chartName string) string {
	return "helm/" + chartName + ".tgz"
}

// GetKubeAPI returns the CoreV1Interface
func GetKubeAPI(useInClusterConfig bool) (v1.CoreV1Interface, error) {

	clientset, err := GetClientset(useInClusterConfig)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
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
		result, getErr := deploymentsClient.Get(context.TODO(), deployment, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}

		result.Spec.Replicas = int32Ptr(replicas)
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	return retryErr
}

func int32Ptr(i int32) *int32 { return &i }

// RestartPodsWithSelector restarts the pods which are found in the provided namespace and selector
func RestartPodsWithSelector(useInClusterConfig bool, namespace string, selector string) error {
	clientset, err := GetKubeAPI(useInClusterConfig)
	if err != nil {
		return err
	}
	pods, err := clientset.Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if err := clientset.Pods(namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
