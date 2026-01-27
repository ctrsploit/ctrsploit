package secret

import (
	"context"
	"fmt"
	"slices"

	"github.com/ctrsploit/sploit-spec/pkg/log"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecretAccessVerbs represents the verbs that allow secret access
var SecretAccessVerbs = []string{"get", "list", "watch", "*"}

// RoleWithSecretAccess represents a role or clusterrole that has access to secrets
type RoleWithSecretAccess struct {
	Kind      string   // "ClusterRole" or "Role"
	Namespace string   // empty for ClusterRole
	Name      string   // role name
	Verbs     []string // verbs that grant secret access
}

// GetRolesWithSecretAccess returns a list of roles/clusterroles that have access to secrets.
func GetRolesWithSecretAccess(clientset kubernetes.Interface) ([]RoleWithSecretAccess, error) {
	var result []RoleWithSecretAccess

	// Get all ClusterRoles
	clusterRoles, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list cluster roles: %w", err)
	}

	for _, role := range clusterRoles.Items {
		if r := buildRoleWithSecretAccess(role.Name, role.Rules, "ClusterRole", ""); r != nil {
			result = append(result, *r)
		}
	}

	// Get all Roles from all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, ns := range namespaces.Items {
		roles, err := clientset.RbacV1().Roles(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Logger.Warnf("Failed to list roles in namespace %s: %v, continuing with other namespaces", ns.Name, err)
			continue
		}

		for _, role := range roles.Items {
			if r := buildRoleWithSecretAccess(role.Name, role.Rules, "Role", role.Namespace); r != nil {
				result = append(result, *r)
			}
		}
	}

	return result, nil
}

// buildRoleWithSecretAccess builds a RoleWithSecretAccess if the role has secret access
func buildRoleWithSecretAccess(name string, rules []rbacv1.PolicyRule, kind, namespace string) *RoleWithSecretAccess {
	verbs := extractSecretAccessVerbs(rules)
	if len(verbs) == 0 {
		return nil
	}
	return &RoleWithSecretAccess{
		Kind:      kind,
		Namespace: namespace,
		Name:      name,
		Verbs:     verbs,
	}
}

// extractSecretAccessVerbs extracts verbs from rules that grant access to secrets.
// It first checks if any rule allows access to "secrets" resource (or "*") with verbs
// that include "get", "list", "watch", or "*". If such a rule exists, it extracts
// all verbs from all rules that apply to secrets (not just the matching verbs).
func extractSecretAccessVerbs(rules []rbacv1.PolicyRule) []string {
	// First, check if any rule has both secrets resource access and matching verbs
	hasQualifyingRule := false
	for _, rule := range rules {
		if !hasSecretResource(rule) {
			continue
		}

		// Check if any verb matches the secret access verbs
		for _, verb := range rule.Verbs {
			if slices.Contains(SecretAccessVerbs, verb) {
				hasQualifyingRule = true
				break
			}
		}

		if hasQualifyingRule {
			break
		}
	}

	// If no qualifying rule found, return empty
	if !hasQualifyingRule {
		return nil
	}

	// Extract all verbs from all rules that apply to secrets
	var verbs []string
	verbSet := make(map[string]bool)

	for _, rule := range rules {
		if !hasSecretResource(rule) {
			continue
		}

		// Extract all verbs from this rule (not just the matching ones)
		for _, verb := range rule.Verbs {
			if !verbSet[verb] {
				verbs = append(verbs, verb)
				verbSet[verb] = true
			}
		}
	}

	return verbs
}

// hasSecretResource checks if a rule applies to secrets resource
func hasSecretResource(rule rbacv1.PolicyRule) bool {
	for _, resource := range rule.Resources {
		if resource == "secrets" || resource == "*" {
			return true
		}
	}
	return false
}
