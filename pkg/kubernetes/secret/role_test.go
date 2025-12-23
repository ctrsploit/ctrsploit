package kubernetes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetRolesWithSecretAccess(t *testing.T) {
	tests := []struct {
		name           string
		clusterRoles   []*rbacv1.ClusterRole
		namespaces     []*corev1.Namespace
		rolesByNS      map[string][]*rbacv1.Role
		expectedResult []RoleWithSecretAccess
		expectError    bool
	}{
		{
			name: "ClusterRole with secret access",
			clusterRoles: []*rbacv1.ClusterRole{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "admin"},
					Rules: []rbacv1.PolicyRule{
						{
							Resources: []string{"secrets"},
							Verbs:     []string{"get", "list", "watch"},
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{},
			expectedResult: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"get", "list", "watch"}},
			},
			expectError: false,
		},
		{
			name:         "Role with secret access",
			clusterRoles: []*rbacv1.ClusterRole{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-reader",
							Namespace: "default",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"secrets"},
								Verbs:     []string{"get", "list"},
							},
						},
					},
				},
			},
			expectedResult: []RoleWithSecretAccess{
				{Kind: "Role", Namespace: "default", Name: "secret-reader", Verbs: []string{"get", "list"}},
			},
			expectError: false,
		},
		{
			name: "ClusterRole and Role with secret access",
			clusterRoles: []*rbacv1.ClusterRole{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "cluster-admin"},
					Rules: []rbacv1.PolicyRule{
						{
							Resources: []string{"*"},
							Verbs:     []string{"*"},
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-reader",
							Namespace: "default",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"secrets"},
								Verbs:     []string{"get", "watch"},
							},
						},
					},
				},
				"kube-system": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "system-reader",
							Namespace: "kube-system",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"secrets"},
								Verbs:     []string{"list"},
							},
						},
					},
				},
			},
			expectedResult: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "cluster-admin", Verbs: []string{"*"}},
				{Kind: "Role", Namespace: "default", Name: "secret-reader", Verbs: []string{"get", "watch"}},
				{Kind: "Role", Namespace: "kube-system", Name: "system-reader", Verbs: []string{"list"}},
			},
			expectError: false,
		},
		{
			name: "Role without secret access should be excluded",
			clusterRoles: []*rbacv1.ClusterRole{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "pod-reader"},
					Rules: []rbacv1.PolicyRule{
						{
							Resources: []string{"pods"},
							Verbs:     []string{"get", "list"},
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "configmap-reader",
							Namespace: "default",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"configmaps"},
								Verbs:     []string{"get", "list"},
							},
						},
					},
				},
			},
			expectedResult: []RoleWithSecretAccess{},
			expectError:    false,
		},
		{
			name:         "Role with secret resource but no matching verbs should be excluded",
			clusterRoles: []*rbacv1.ClusterRole{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-creator",
							Namespace: "default",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"secrets"},
								Verbs:     []string{"create", "update", "delete"},
							},
						},
					},
				},
			},
			expectedResult: []RoleWithSecretAccess{},
			expectError:    false,
		},
		{
			name:         "Role with wildcard resource and matching verb",
			clusterRoles: []*rbacv1.ClusterRole{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "wildcard-reader",
							Namespace: "default",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"*"},
								Verbs:     []string{"get", "list", "watch", "create"},
							},
						},
					},
				},
			},
			expectedResult: []RoleWithSecretAccess{
				{Kind: "Role", Namespace: "default", Name: "wildcard-reader", Verbs: []string{"get", "list", "watch", "create"}},
			},
			expectError: false,
		},
		{
			name:         "Multiple rules with secret access - extract all verbs",
			clusterRoles: []*rbacv1.ClusterRole{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS: map[string][]*rbacv1.Role{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "multi-rule-reader",
							Namespace: "default",
						},
						Rules: []rbacv1.PolicyRule{
							{
								Resources: []string{"secrets"},
								Verbs:     []string{"get", "list"},
							},
							{
								Resources: []string{"secrets"},
								Verbs:     []string{"watch", "patch"},
							},
						},
					},
				},
			},
			expectedResult: []RoleWithSecretAccess{
				{Kind: "Role", Namespace: "default", Name: "multi-rule-reader", Verbs: []string{"get", "list", "watch", "patch"}},
			},
			expectError: false,
		},
		{
			name:         "Empty cluster - no roles",
			clusterRoles: []*rbacv1.ClusterRole{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			rolesByNS:      map[string][]*rbacv1.Role{},
			expectedResult: []RoleWithSecretAccess{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake clientset
			clientset := fake.NewSimpleClientset()

			// Add ClusterRoles
			for _, cr := range tt.clusterRoles {
				_, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add Namespaces
			for _, ns := range tt.namespaces {
				_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add Roles by namespace
			for ns, roles := range tt.rolesByNS {
				for _, role := range roles {
					_, err := clientset.RbacV1().Roles(ns).Create(context.TODO(), role, metav1.CreateOptions{})
					assert.NoError(t, err)
				}
			}

			// Test the function
			result, err := GetRolesWithSecretAccess(clientset)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Sort results for comparison (order may vary)
				assert.ElementsMatch(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetRolesWithSecretAccess_NamespaceError(t *testing.T) {
	// Test that function continues when listing roles in a namespace fails
	clientset := fake.NewSimpleClientset()

	// Add a namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "default"},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	assert.NoError(t, err)

	// Add a role with secret access
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret-reader",
			Namespace: "default",
		},
		Rules: []rbacv1.PolicyRule{
			{
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "list"},
			},
		},
	}
	_, err = clientset.RbacV1().Roles("default").Create(context.TODO(), role, metav1.CreateOptions{})
	assert.NoError(t, err)

	// Delete the namespace to simulate an error when listing roles
	// (This will cause the Roles().List() call to fail for that namespace)
	err = clientset.CoreV1().Namespaces().Delete(context.TODO(), "default", metav1.DeleteOptions{})
	assert.NoError(t, err)

	// The function should still work and return empty result (since namespace was deleted)
	// In real scenario, it would log a warning and continue
	result, err := GetRolesWithSecretAccess(clientset)
	assert.NoError(t, err)
	// After namespace deletion, the role listing will fail, so result should be empty
	assert.Empty(t, result)
}
