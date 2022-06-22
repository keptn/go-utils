package kubeutils

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"

	"github.com/keptn/go-utils/pkg/common/fileutils"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientset returns the kubernetes Clientset
func GetClientSet(useInClusterConfig bool) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		var kubeconfig string
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = fileutils.ExpandTilde(os.Getenv("KUBECONFIG"))
		} else {
			kubeconfig = filepath.Join(
				fileutils.UserHomeDir(), ".kube", "config",
			)
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
