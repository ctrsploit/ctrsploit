package kubernetes

import (
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesClient(kubeconfigPath ...string) (*kubernetes.Clientset, error) {
	config, err := GetKubernetesConfig(kubeconfigPath...)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %v", err)
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return c, nil
}

func GetKubernetesConfig(kubeconfigPath ...string) (*rest.Config, error) {
	if len(kubeconfigPath) > 0 && kubeconfigPath[0] != "" {
		log.Logger.Infof("Using kubeconfig from path: %s", kubeconfigPath[0])
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath[0])
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		log.Logger.Debugf("Fall back to in-cluster config")
		return rest.InClusterConfig()
	} else {
		log.Logger.Debugf("Using kubeconfig from default location")
	}

	return config, nil
}
