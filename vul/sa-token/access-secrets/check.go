package access_secrets

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes/secret"
	prerequisiteSecret "github.com/ctrsploit/ctrsploit/prerequisite/kubernetes/service-account/secret"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

// Check checks for pods with secret access and prints the analysis results.
// This function can be used by other projects to check and report secret access.
func Check() error {
	pods, err := prerequisiteSecret.GetPodsWithSecretAccess(&prerequisiteSecret.HasPodsWithSecretAccess)
	if err != nil {
		return fmt.Errorf("failed to get pods with secret access: %w", err)
	}
	Report(pods)
	return nil
}

// Report prints the analysis results for the given pods with secret access.
// This function can be used by other projects to print secret access analysis results.
func Report(pods []secret.PodWithSecretAccess) {
	if len(pods) == 0 {
		log.Logger.Infof("No pods with secret access found")
		return
	}

	// Collect cluster-wide statistics
	allNamespaces := make(map[string]bool)
	allPermissions := make(map[string]bool)
	namespacePodCount := make(map[string]int)
	serviceAccountMap := make(map[string]map[string]bool) // namespace -> serviceAccount -> true
	namespacePermissions := make(map[string][]secret.ServiceAccountPermission)

	for _, pod := range pods {
		allNamespaces[pod.Namespace] = true
		namespacePodCount[pod.Namespace]++

		// Track service accounts
		if serviceAccountMap[pod.Namespace] == nil {
			serviceAccountMap[pod.Namespace] = make(map[string]bool)
		}
		serviceAccountMap[pod.Namespace][pod.ServiceAccountName] = true

		// Track unique permissions
		for _, perm := range pod.Permissions {
			permKey := fmt.Sprintf("%s/%s/%s", perm.BindingKind, perm.RoleRefKind, perm.RoleRefName)
			allPermissions[permKey] = true

			// Track permissions per namespace (deduplicate)
			seenPermKey := fmt.Sprintf("%s/%s/%s/%s", perm.BindingKind, perm.RoleRefKind, perm.RoleRefName, strings.Join(perm.Verbs, ","))
			if namespacePermissions[pod.Namespace] == nil {
				namespacePermissions[pod.Namespace] = make([]secret.ServiceAccountPermission, 0)
			}
			// Check if this permission is already in the list for this namespace
			found := false
			for _, existingPerm := range namespacePermissions[pod.Namespace] {
				existingPermKey := fmt.Sprintf("%s/%s/%s/%s", existingPerm.BindingKind, existingPerm.RoleRefKind, existingPerm.RoleRefName, strings.Join(existingPerm.Verbs, ","))
				if existingPermKey == seenPermKey {
					found = true
					break
				}
			}
			if !found {
				namespacePermissions[pod.Namespace] = append(namespacePermissions[pod.Namespace], perm)
			}
		}
	}

	// Display cluster-wide information
	log.Logger.Infof("")
	log.Logger.Infof("=== Cluster-wide Pods with Secret Access ===")
	log.Logger.Infof("Found %d pod(s) with secret access permissions across the cluster", len(pods))

	// Display pods grouped by namespace
	log.Logger.Infof("")
	log.Logger.Infof("=== Pods by Namespace ===")
	for namespace := range allNamespaces {
		log.Logger.Infof("")
		log.Logger.Infof("Namespace: %s", namespace)
		log.Logger.Infof("  Pods: %d", namespacePodCount[namespace])
		log.Logger.Infof("  Service Accounts: %d", len(serviceAccountMap[namespace]))

		// Show unique RBAC permissions for this namespace
		if perms, exists := namespacePermissions[namespace]; exists && len(perms) > 0 {
			log.Logger.Infof("  RBAC Permissions:")
			for _, perm := range perms {
				log.Logger.Infof("    - [%s] %s (Scope: %s, Verbs: %s)", perm.BindingKind, perm.RoleRefName, perm.Scope(), strings.Join(perm.Verbs, ","))
			}
		}
	}

	// Display summary
	log.Logger.Infof("")
	log.Logger.Infof("=== Summary ===")
	log.Logger.Infof("Total pods with secret access: %d", len(pods))
	log.Logger.Infof("Namespaces affected: %d", len(allNamespaces))
	log.Logger.Infof("Unique RBAC permissions: %d", len(allPermissions))
}
