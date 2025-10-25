package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesClient(kubeconfigPath ...string) (*kubernetes.Clientset, error) {
	config, err := GetKubernetesConfig(kubeconfigPath...)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return clientset, nil
}

func GetKubernetesConfig(kubeconfigPath ...string) (*rest.Config, error) {
	if len(kubeconfigPath) > 0 && kubeconfigPath[0] != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath[0])
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return rest.InClusterConfig()
	}

	return config, nil
}
