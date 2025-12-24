package secret

import (
	"context"
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodWithSecretAccess represents a Pod and its associated secret access permissions
// through the ServiceAccount it uses.
type PodWithSecretAccess struct {
	// Namespace is the namespace of the Pod
	Namespace string `json:"namespace"`
	// Name is the name of the Pod
	Name string `json:"name"`
	// ServiceAccountName is the name of the ServiceAccount used by the Pod
	ServiceAccountName string `json:"serviceAccountName"`
	// Permissions is the list of permissions granted to the ServiceAccount used by this Pod
	Permissions []ServiceAccountPermission `json:"permissions"`
}

// GetPodsWithSecretAccess scans all Pods and returns those that use ServiceAccounts
// with secret access permissions. It matches Pods with their ServiceAccount permissions
// based on the ServiceAccountWithPermissions list provided.
//
// 1. Gets all Pods across all namespaces
// 2. For each Pod, checks if its ServiceAccount has secret access permissions
// 3. Returns Pod information along with the associated permissions
func GetPodsWithSecretAccess(clientset kubernetes.Interface, saWithPerms []ServiceAccountWithPermissions) ([]PodWithSecretAccess, error) {
	// Build a map for quick lookup: ServiceAccount key (namespace/name) -> []ServiceAccountPermission
	saMap := make(map[string][]ServiceAccountPermission)
	for _, sa := range saWithPerms {
		saKey := buildServiceAccountKey(sa.Namespace, sa.Name)
		saMap[saKey] = sa.Permissions
	}

	result := make([]PodWithSecretAccess, 0)

	// Get all Pods from all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, ns := range namespaces.Items {
		pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Logger.Warnf("Failed to list pods in namespace %s: %v, continuing with other namespaces", ns.Name, err)
			continue
		}

		for _, pod := range pods.Items {
			// Get the ServiceAccount name used by the Pod
			// If spec.serviceAccountName is empty, it defaults to "default"
			saName := pod.Spec.ServiceAccountName
			if saName == "" {
				saName = "default"
			}

			// Build the ServiceAccount key
			saKey := buildServiceAccountKey(pod.Namespace, saName)

			// Check if this ServiceAccount has secret access permissions
			if perms, exists := saMap[saKey]; exists && len(perms) > 0 {
				// For each permission, create a PodWithSecretAccess entry
				// This matches the jq logic that outputs one line per permission
				for _, perm := range perms {
					result = append(result, PodWithSecretAccess{
						Namespace:          pod.Namespace,
						Name:               pod.Name,
						ServiceAccountName: saName,
						Permissions:        []ServiceAccountPermission{perm},
					})
				}
			}
		}
	}

	return result, nil
}

// GetAllPodsWithSecretAccess is a convenience function that chains the three analysis functions together:
// 1. GetRolesWithSecretAccess - finds all Roles/ClusterRoles with secret access
// 2. GetServiceAccountsWithSecretAccess - finds all ServiceAccounts bound to those roles
// 3. GetPodsWithSecretAccess - finds all Pods using those ServiceAccounts
//
// This function provides a complete analysis chain from RBAC roles to actual Pods that have
// secret access permissions through their ServiceAccounts.
func GetAllPodsWithSecretAccess(clientset kubernetes.Interface) ([]PodWithSecretAccess, error) {
	// Step 1: Get all roles/clusterroles with secret access
	roles, err := GetRolesWithSecretAccess(clientset)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles with secret access: %w", err)
	}

	// Step 2: Get all service accounts with secret access permissions
	serviceAccounts, err := GetServiceAccountsWithSecretAccess(clientset, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to get service accounts with secret access: %w", err)
	}

	// Step 3: Get all pods using those service accounts
	pods, err := GetPodsWithSecretAccess(clientset, serviceAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods with secret access: %w", err)
	}

	return pods, nil
}
