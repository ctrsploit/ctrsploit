package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ServiceAccountPermission represents permission information for a ServiceAccount
// that has access to secrets through RBAC bindings.
// The field names follow Kubernetes RBAC standard naming conventions.
type ServiceAccountPermission struct {
	// RoleRefKind is the kind of role: "Role" or "ClusterRole"
	RoleRefKind string `json:"roleRefKind"`
	// RoleRefName is the name of the role or cluster role
	RoleRefName string `json:"roleRefName"`
	// BindingKind is the kind of binding: "RoleBinding" or "ClusterRoleBinding"
	BindingKind string `json:"bindingKind"`
	// Namespace is the namespace where the binding applies.
	// Empty string for ClusterRoleBinding, or the namespace for RoleBinding
	Namespace string `json:"namespace,omitempty"`
	// Verbs are the RBAC verbs that grant secret access (e.g., ["get", "list", "watch", "*"])
	Verbs []string `json:"verbs"`
}

// String returns a formatted string representation of the permission.
// Format: "[RoleRefKind] RoleRefName (Scope: Namespace, Verbs: verb1,verb2,...)"
func (p ServiceAccountPermission) String() string {
	scope := p.Scope()
	verbsStr := strings.Join(p.Verbs, ",")
	return fmt.Sprintf("[%s] %s (Scope: %s, Verbs: %s)", p.RoleRefKind, p.RoleRefName, scope, verbsStr)
}

// HasVerb checks if the permission includes the specified verb.
// Returns true if verbs contain the verb or "*".
func (p ServiceAccountPermission) HasVerb(verb string) bool {
	for _, v := range p.Verbs {
		if v == verb || v == "*" {
			return true
		}
	}
	return false
}

// IsClusterWide returns true if this is a cluster-wide permission (ClusterRoleBinding).
func (p ServiceAccountPermission) IsClusterWide() bool {
	return p.BindingKind == "ClusterRoleBinding"
}

// Scope returns a human-readable scope description.
func (p ServiceAccountPermission) Scope() string {
	if p.IsClusterWide() {
		return "Cluster-Wide"
	}
	if p.Namespace != "" {
		return fmt.Sprintf("NS: %s", p.Namespace)
	}
	return "Unknown"
}

// RoleDisplayName returns a formatted role name for display.
// Format: "[RoleRefKind] RoleRefName"
func (p ServiceAccountPermission) RoleDisplayName() string {
	return fmt.Sprintf("[%s] %s", p.RoleRefKind, p.RoleRefName)
}

// ServiceAccountWithPermissions represents a ServiceAccount and its associated permissions.
type ServiceAccountWithPermissions struct {
	// Namespace is the namespace of the ServiceAccount
	Namespace string `json:"namespace"`
	// Name is the name of the ServiceAccount
	Name string `json:"name"`
	// Permissions is the list of permissions granted to this ServiceAccount
	Permissions []ServiceAccountPermission `json:"permissions"`
}

// GetServiceAccountsWithSecretAccess returns a list of ServiceAccounts with their secret access permissions.
// It scans all ClusterRoleBindings and RoleBindings to find ServiceAccounts that are bound to roles with secret access.
func GetServiceAccountsWithSecretAccess(clientset kubernetes.Interface, roles []RoleWithSecretAccess) ([]ServiceAccountWithPermissions, error) {
	// Build a map for quick lookup: roleKey -> RoleWithSecretAccess
	roleMap := make(map[string]*RoleWithSecretAccess)
	for i := range roles {
		role := &roles[i]
		roleKey := buildRoleKey(role.Kind, role.Namespace, role.Name)
		roleMap[roleKey] = role
	}

	// Result map: ServiceAccount key (namespace/name) -> []ServiceAccountPermission
	// We use a map temporarily to aggregate permissions, then convert to list
	resultMap := make(map[string][]ServiceAccountPermission)

	// Get all ClusterRoleBindings
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list cluster role bindings: %w", err)
	}

	for _, binding := range clusterRoleBindings.Items {
		roleKey := buildRoleKey("ClusterRole", "", binding.RoleRef.Name)
		if role, exists := roleMap[roleKey]; exists {
			permission := ServiceAccountPermission{
				RoleRefKind: binding.RoleRef.Kind,
				RoleRefName: binding.RoleRef.Name,
				BindingKind: "ClusterRoleBinding",
				Namespace:   "",
				Verbs:       role.Verbs,
			}
			addServiceAccountPermissions(binding.Subjects, permission, resultMap)
		}
	}

	// Get all RoleBindings from all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, ns := range namespaces.Items {
		roleBindings, err := clientset.RbacV1().RoleBindings(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Logger.Warnf("Failed to list role bindings in namespace %s: %v, continuing with other namespaces", ns.Name, err)
			continue
		}

		for _, binding := range roleBindings.Items {
			// For RoleBinding, check if it references a Role or ClusterRole
			switch binding.RoleRef.Kind {
			case "Role":
				roleKey := buildRoleKey("Role", binding.Namespace, binding.RoleRef.Name)
				if role, exists := roleMap[roleKey]; exists {
					permission := ServiceAccountPermission{
						RoleRefKind: binding.RoleRef.Kind,
						RoleRefName: binding.RoleRef.Name,
						BindingKind: "RoleBinding",
						Namespace:   binding.Namespace,
						Verbs:       role.Verbs,
					}
					addServiceAccountPermissions(binding.Subjects, permission, resultMap)
				}
			case "ClusterRole":
				// RoleBinding can also reference a ClusterRole
				roleKey := buildRoleKey("ClusterRole", "", binding.RoleRef.Name)
				if role, exists := roleMap[roleKey]; exists {
					permission := ServiceAccountPermission{
						RoleRefKind: binding.RoleRef.Kind,
						RoleRefName: binding.RoleRef.Name,
						BindingKind: "RoleBinding",
						Namespace:   binding.Namespace,
						Verbs:       role.Verbs,
					}
					addServiceAccountPermissions(binding.Subjects, permission, resultMap)
				}
			}
		}
	}

	// Convert map to list of ServiceAccountWithPermissions
	result := make([]ServiceAccountWithPermissions, 0, len(resultMap))
	for key, perms := range resultMap {
		parts := strings.Split(key, "/")
		if len(parts) == 2 {
			result = append(result, ServiceAccountWithPermissions{
				Namespace:   parts[0],
				Name:        parts[1],
				Permissions: perms,
			})
		}
	}

	return result, nil
}

// buildRoleKey builds a key for role lookup in the format: "kind;namespace;name"
// For ClusterRole, namespace is empty
func buildRoleKey(kind, namespace, name string) string {
	return fmt.Sprintf("%s;%s;%s", kind, namespace, name)
}

// buildServiceAccountKey builds a key for ServiceAccount in the format: "namespace/name"
func buildServiceAccountKey(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

// addServiceAccountPermissions extracts ServiceAccount subjects from the binding and adds them to the result map.
// It filters subjects to only include those with Kind "ServiceAccount" and adds the permission to each matching ServiceAccount.
func addServiceAccountPermissions(subjects []rbacv1.Subject, permission ServiceAccountPermission, resultMap map[string][]ServiceAccountPermission) {
	for _, subject := range subjects {
		if subject.Kind == "ServiceAccount" {
			saKey := buildServiceAccountKey(subject.Namespace, subject.Name)
			resultMap[saKey] = append(resultMap[saKey], permission)
		}
	}
}
