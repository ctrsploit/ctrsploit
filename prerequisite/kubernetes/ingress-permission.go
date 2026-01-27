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
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressPermission struct {
	prerequisite.BasePrerequisite
	Action string
}

var HasIngressCreatePermission = IngressPermission{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "HasIngressCreatePermission",
		Info:   "Check if current user has ingress creation permission",
		ExeEnv: exeenv.K8S,
	},
	Action: "create",
}

var HasIngressUpdatePermission = IngressPermission{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "HasIngressUpdatePermission",
		Info:   "Check if current user has ingress update permission",
		ExeEnv: exeenv.K8S,
	},
	Action: "update",
}

func (p *IngressPermission) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	log.Logger.Debugf("Checking IngressPermission (Action: %s)", p.Action)
	clientset, err := kubernetes.
		GetKubernetesClient()
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
		log.Logger.Debugf("Checking ingress creation permission")
		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name: "permission-test-ingress",
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: "test.example.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     "/test",
										PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: "test-service",
												Port: networkingv1.ServiceBackendPort{
													Number: 80,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = clientset.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
		if err != nil {
			err = fmt.Errorf("no ingress create permission: %v", err)
			awesome_error.CheckErr(err)
			return false, err
		}
		log.Logger.Debugf("Ingress creation permission check passed")

	case "update":
		log.Logger.Debugf("Checking ingress update permission")
		_, err = clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			err = fmt.Errorf("no ingress list/update permission: %v", err)
			awesome_error.CheckErr(err)
			return false, err
		}
		log.Logger.Debugf("Ingress update permission check passed")
	}

	p.Satisfied = true
	log.Logger.Debugf("IngressPermission.Satisfied = %v", p.Satisfied)
	p.Checked = true
	return p.Satisfied, nil
}
