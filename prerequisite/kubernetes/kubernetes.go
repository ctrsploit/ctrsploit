package kubernetes

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

func getKubernetesClient(kubeconfigPath ...string) (*k8sclient.Clientset, error) {
	return kubernetes.GetKubernetesClient(kubeconfigPath...)
}

type KubeletVersionInfo struct {
	NodeName string
	Version  string
}

func GetKubeletVersions() ([]KubeletVersionInfo, error) {
	clientset, err := getKubernetesClient()
	if err != nil {
		log.Logger.Errorf("Failed to get Kubernetes client: %v", err)
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}

	log.Logger.Debugf("Found %d nodes", len(nodes.Items))

	var versions []KubeletVersionInfo
	for _, node := range nodes.Items {
		kubeletVersion := node.Status.NodeInfo.KubeletVersion
		log.Logger.Debugf("Node %s kubelet version: %s", node.Name, kubeletVersion)
		versions = append(versions, KubeletVersionInfo{
			NodeName: node.Name,
			Version:  kubeletVersion,
		})
	}

	return versions, nil
}

func GetVulnerableNodesToCVE202125741() ([]string, error) {
	versions, err := GetKubeletVersions()
	if err != nil {
		return nil, err
	}

	var vulnerableNodes []string
	for _, versionInfo := range versions {
		if isVulnerableToCVE202125741(versionInfo.Version) {
			vulnerableNodes = append(vulnerableNodes, versionInfo.NodeName)
			log.Logger.Debugf("Node %s is vulnerable to CVE-2021-25741 (version: %s)", versionInfo.NodeName, versionInfo.Version)
		}
	}

	return vulnerableNodes, nil
}

type CVE202125741KubeletVersion struct {
	prerequisite.BasePrerequisite
	VulnerableNodes []string
}

var VulnerableToCVE202125741 = CVE202125741KubeletVersion{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "VulnerableToCVE202125741",
		Info:   "Check if kubelet version is vulnerable to CVE-2021-25741",
		ExeEnv: exeenv.Local | exeenv.K8S,
	},
	VulnerableNodes: []string{},
}

func (p *CVE202125741KubeletVersion) Check() error {
	log.Logger.Debugf("Checking CVE-2021-25741 KubeletVersion")
	
	err := p.BasePrerequisite.Check()
	if err != nil {
		log.Logger.Errorf("BasePrerequisite.Check() failed: %v", err)
		return err
	}

	vulnerableNodes, err := GetVulnerableNodesToCVE202125741()
	if err != nil {
		return err
	}

	if len(vulnerableNodes) == 0 {
		return fmt.Errorf("not found kubelet version vulnerable to CVE-2021-25741")
	}

	p.VulnerableNodes = vulnerableNodes
	p.Satisfied = true
	log.Logger.Infof("Found %d vulnerable nodes: %v", len(vulnerableNodes), vulnerableNodes)
	return nil
}

func (p *CVE202125741KubeletVersion) GetVulnerableNodes() []string {
	return p.VulnerableNodes
}

func isVulnerableToCVE202125741(version string) bool {
	version = strings.TrimPrefix(version, "v")
	
	// CVE-2021-25741 affected versions:
	// v1.22.0 - v1.22.1
	// v1.21.0 - v1.21.4  
	// v1.20.0 - v1.20.10
	// <= v1.19.14
	
	if strings.HasPrefix(version, "1.22.") {
		return version == "1.22.0" || version == "1.22.1"
	}
	
	if strings.HasPrefix(version, "1.21.") {
		minorVersion := strings.TrimPrefix(version, "1.21.")
		return compareVersion(minorVersion, "4") <= 0
	}
	
	if strings.HasPrefix(version, "1.20.") {
		minorVersion := strings.TrimPrefix(version, "1.20.")
		return compareVersion(minorVersion, "10") <= 0
	}
	
	if strings.HasPrefix(version, "1.19.") {
		minorVersion := strings.TrimPrefix(version, "1.19.")
		return compareVersion(minorVersion, "14") <= 0
	}
	
	if strings.HasPrefix(version, "1.18.") || 
	   strings.HasPrefix(version, "1.17.") ||
	   strings.HasPrefix(version, "1.16.") ||
	   strings.HasPrefix(version, "1.15.") {
		return true
	}
	
	return false
}

func compareVersion(version1 string, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")
	
	for i := 0; i < len(v1) && i < len(v2); i++ {
		if v1[i] < v2[i] {
			return -1
		} else if v1[i] > v2[i] {
			return 1
		}
	}
	
	if len(v1) < len(v2) {
		return -1
	} else if len(v1) > len(v2) {
		return 1
	}
	
	return 0
}

type PodPermission struct {
	prerequisite.BasePrerequisite
	Action string
}

var HasPodCreatePermission = PodPermission{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "HasPodCreatePermission",
		Info:   "Check if current user has pod creation permission",
		ExeEnv: exeenv.Local | exeenv.K8S,
	},
	Action: "create",
}

var HasPodExecPermission = PodPermission{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "HasPodExecPermission", 
		Info:   "Check if current user has pod exec permission",
		ExeEnv: exeenv.Local | exeenv.K8S,
	},
	Action: "exec",
}

func (p *PodPermission) Check() error {
	log.Logger.Debugf("Checking PodPermission (Action: %s)", p.Action)
	
	err := p.BasePrerequisite.Check()
	if err != nil {
		log.Logger.Errorf("BasePrerequisite.Check() failed: %v", err)
		return err
	}

	clientset, err := getKubernetesClient()
	if err != nil {
		log.Logger.Errorf("Failed to get Kubernetes client: %v", err)
		return err
	}

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	log.Logger.Debugf("Checking namespace: %s", namespace)

	switch p.Action {
	case "create":
		log.Logger.Debugf("Checking pod creation permission")
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "permission-test-pod",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "test",
						Image: "alpine:latest",
					},
				},
			},
		}
		
		_, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
		if err != nil {
			return fmt.Errorf("no pod create permission: %v", err)
		}
		log.Logger.Debugf("Pod creation permission check passed")
		
	case "exec":
		log.Logger.Debugf("Checking pod exec permission")
		_, err = clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("no pod list/exec permission: %v", err)
		}
		log.Logger.Debugf("Pod exec permission check passed")
	}

	p.Satisfied = true
	log.Logger.Debugf("PodPermission.Satisfied = %v", p.Satisfied)
	return nil
} 