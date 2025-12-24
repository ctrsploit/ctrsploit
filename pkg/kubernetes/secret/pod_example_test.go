package secret

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// ExampleGetAllPodsWithSecretAccess demonstrates how to use GetAllPodsWithSecretAccess
// to find all Pods that have secret access permissions through their ServiceAccounts.
func ExampleGetAllPodsWithSecretAccess() {
	// Create a fake Kubernetes clientset for demonstration
	// In real usage, you would use:
	//   clientset, err := GetKubernetesClient()
	clientset := fake.NewSimpleClientset()

	// Setup: Create a ClusterRole with secret access
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{Name: "secret-reader"},
		Rules: []rbacv1.PolicyRule{
			{
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
	clientset.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})

	// Setup: Create a ClusterRoleBinding that binds the role to a ServiceAccount
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: "secret-reader-binding"},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "secret-reader",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Namespace: "default",
				Name:      "my-sa",
			},
		},
	}
	clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})

	// Setup: Create a namespace
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "default"},
	}
	clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})

	// Setup: Create a Pod that uses the ServiceAccount
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "my-sa",
		},
	}
	clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})

	// Use GetAllPodsWithSecretAccess to find all Pods with secret access
	pods, err := GetAllPodsWithSecretAccess(clientset)
	if err != nil {
		log.Fatalf("Failed to get pods with secret access: %v", err)
	}

	// Print results
	fmt.Printf("Found %d pod(s) with secret access:\n", len(pods))
	for _, pod := range pods {
		fmt.Printf("  Pod: %s/%s (ServiceAccount: %s)\n", pod.Namespace, pod.Name, pod.ServiceAccountName)
		for _, perm := range pod.Permissions {
			fmt.Printf("    Permission: [%s] %s via %s (Verbs: %v)\n",
				perm.RoleRefKind, perm.RoleRefName, perm.BindingKind, perm.Verbs)
		}
	}

	// Output:
	// Found 1 pod(s) with secret access:
	//   Pod: default/my-pod (ServiceAccount: my-sa)
	//     Permission: [ClusterRole] secret-reader via ClusterRoleBinding (Verbs: [get list watch])
}
