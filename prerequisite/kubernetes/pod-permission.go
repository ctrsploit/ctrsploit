package kubernetes

import (
	"context"
	"fmt"
	"os"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func (p *PodPermission) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	log.Logger.Debugf("Checking PodPermission (Action: %s)", p.Action)
	clientset, err := kubernetes.GetKubernetesClient()
	if err != nil {
		return false, fmt.Errorf("failed to get Kubernetes client: %v", err)
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
			err = fmt.Errorf("no pod create permission: %v", err)
			awesome_error.CheckErr(err)
			return false, err
		}
		log.Logger.Debugf("Pod creation permission check passed")

	case "exec":
		log.Logger.Debugf("Checking pod exec permission")
		_, err = clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			err = fmt.Errorf("no pod list/exec permission: %v", err)
			awesome_error.CheckErr(err)
			return false, err
		}
		log.Logger.Debugf("Pod exec permission check passed")
	}

	p.Satisfied = true
	log.Logger.Debugf("PodPermission.Satisfied = %v", p.Satisfied)
	p.Checked = true
	return p.Satisfied, nil
}
