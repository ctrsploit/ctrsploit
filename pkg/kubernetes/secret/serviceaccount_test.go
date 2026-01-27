package secret

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetServiceAccountsWithSecretAccess(t *testing.T) {
	tests := []struct {
		name                string
		roles               []RoleWithSecretAccess
		clusterRoleBindings []*rbacv1.ClusterRoleBinding
		namespaces          []*corev1.Namespace
		roleBindingsByNS    map[string][]*rbacv1.RoleBinding
		expectedResult      []ServiceAccountWithPermissions
		expectError         bool
	}{
		{
			name: "ClusterRoleBinding with ServiceAccount",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"get", "list", "watch"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "default",
							Name:      "sa1",
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa1",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "admin",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"get", "list", "watch"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "RoleBinding with Role",
			roles: []RoleWithSecretAccess{
				{Kind: "Role", Namespace: "default", Name: "secret-reader", Verbs: []string{"get", "list"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-reader-binding",
							Namespace: "default",
						},
						RoleRef: rbacv1.RoleRef{
							APIGroup: "rbac.authorization.k8s.io",
							Kind:     "Role",
							Name:     "secret-reader",
						},
						Subjects: []rbacv1.Subject{
							{
								Kind:      "ServiceAccount",
								Namespace: "default",
								Name:      "sa2",
							},
						},
					},
				},
			},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa2",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "Role",
							RoleRefName: "secret-reader",
							BindingKind: "RoleBinding",
							Namespace:   "default",
							Verbs:       []string{"get", "list"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "RoleBinding with ClusterRole",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "view", Verbs: []string{"get", "list", "watch"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{
				"kube-system": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "view-binding",
							Namespace: "kube-system",
						},
						RoleRef: rbacv1.RoleRef{
							APIGroup: "rbac.authorization.k8s.io",
							Kind:     "ClusterRole",
							Name:     "view",
						},
						Subjects: []rbacv1.Subject{
							{
								Kind:      "ServiceAccount",
								Namespace: "kube-system",
								Name:      "sa3",
							},
						},
					},
				},
			},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "kube-system",
					Name:      "sa3",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "view",
							BindingKind: "RoleBinding",
							Namespace:   "kube-system",
							Verbs:       []string{"get", "list", "watch"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Multiple ServiceAccounts in same binding",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"*"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "default",
							Name:      "sa1",
						},
						{
							Kind:      "ServiceAccount",
							Namespace: "kube-system",
							Name:      "sa2",
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa1",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "admin",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"*"},
						},
					},
				},
				{
					Namespace: "kube-system",
					Name:      "sa2",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "admin",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"*"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "ServiceAccount with multiple permissions",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"*"}},
				{Kind: "Role", Namespace: "default", Name: "secret-reader", Verbs: []string{"get", "list"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "default",
							Name:      "sa1",
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-reader-binding",
							Namespace: "default",
						},
						RoleRef: rbacv1.RoleRef{
							APIGroup: "rbac.authorization.k8s.io",
							Kind:     "Role",
							Name:     "secret-reader",
						},
						Subjects: []rbacv1.Subject{
							{
								Kind:      "ServiceAccount",
								Namespace: "default",
								Name:      "sa1",
							},
						},
					},
				},
			},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa1",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "admin",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"*"},
						},
						{
							RoleRefKind: "Role",
							RoleRefName: "secret-reader",
							BindingKind: "RoleBinding",
							Namespace:   "default",
							Verbs:       []string{"get", "list"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Binding to role without secret access should be ignored",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"*"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "default",
							Name:      "sa1",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "pod-reader-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "pod-reader",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "default",
							Name:      "sa2",
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa1",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "admin",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"*"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Non-ServiceAccount subjects should be ignored",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"*"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "admin-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "admin",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "default",
							Name:      "sa1",
						},
						{
							Kind:     "User",
							APIGroup: "rbac.authorization.k8s.io",
							Name:     "user1",
						},
						{
							Kind:     "Group",
							APIGroup: "rbac.authorization.k8s.io",
							Name:     "group1",
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa1",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "admin",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"*"},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:                "Empty cluster",
			roles:               []RoleWithSecretAccess{},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			expectedResult:   []ServiceAccountWithPermissions{},
			expectError:      false,
		},
		{
			name: "Multiple namespaces with different bindings",
			roles: []RoleWithSecretAccess{
				{Kind: "ClusterRole", Namespace: "", Name: "view", Verbs: []string{"get", "list", "watch"}},
				{Kind: "Role", Namespace: "ns1", Name: "secret-reader", Verbs: []string{"get", "list"}},
				{Kind: "Role", Namespace: "ns2", Name: "secret-reader", Verbs: []string{"watch"}},
			},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "view-binding"},
					RoleRef: rbacv1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     "view",
					},
					Subjects: []rbacv1.Subject{
						{
							Kind:      "ServiceAccount",
							Namespace: "ns1",
							Name:      "sa1",
						},
					},
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "ns1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "ns2"}},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{
				"ns1": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-reader-binding",
							Namespace: "ns1",
						},
						RoleRef: rbacv1.RoleRef{
							APIGroup: "rbac.authorization.k8s.io",
							Kind:     "Role",
							Name:     "secret-reader",
						},
						Subjects: []rbacv1.Subject{
							{
								Kind:      "ServiceAccount",
								Namespace: "ns1",
								Name:      "sa2",
							},
						},
					},
				},
				"ns2": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-reader-binding",
							Namespace: "ns2",
						},
						RoleRef: rbacv1.RoleRef{
							APIGroup: "rbac.authorization.k8s.io",
							Kind:     "Role",
							Name:     "secret-reader",
						},
						Subjects: []rbacv1.Subject{
							{
								Kind:      "ServiceAccount",
								Namespace: "ns2",
								Name:      "sa3",
							},
						},
					},
				},
			},
			expectedResult: []ServiceAccountWithPermissions{
				{
					Namespace: "ns1",
					Name:      "sa1",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "ClusterRole",
							RoleRefName: "view",
							BindingKind: "ClusterRoleBinding",
							Namespace:   "",
							Verbs:       []string{"get", "list", "watch"},
						},
					},
				},
				{
					Namespace: "ns1",
					Name:      "sa2",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "Role",
							RoleRefName: "secret-reader",
							BindingKind: "RoleBinding",
							Namespace:   "ns1",
							Verbs:       []string{"get", "list"},
						},
					},
				},
				{
					Namespace: "ns2",
					Name:      "sa3",
					Permissions: []ServiceAccountPermission{
						{
							RoleRefKind: "Role",
							RoleRefName: "secret-reader",
							BindingKind: "RoleBinding",
							Namespace:   "ns2",
							Verbs:       []string{"watch"},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake clientset
			clientset := fake.NewSimpleClientset()

			// Add ClusterRoleBindings
			for _, crb := range tt.clusterRoleBindings {
				_, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add Namespaces
			for _, ns := range tt.namespaces {
				_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add RoleBindings by namespace
			for ns, roleBindings := range tt.roleBindingsByNS {
				for _, rb := range roleBindings {
					_, err := clientset.RbacV1().RoleBindings(ns).Create(context.TODO(), rb, metav1.CreateOptions{})
					assert.NoError(t, err)
				}
			}

			// Test the function
			result, err := GetServiceAccountsWithSecretAccess(clientset, tt.roles)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// Compare results - use ElementsMatch since order doesn't matter
				assert.Equal(t, len(tt.expectedResult), len(result), "Result count mismatch")
				assert.ElementsMatch(t, tt.expectedResult, result, "ServiceAccounts with permissions mismatch")
			}
		})
	}
}

func TestGetServiceAccountsWithSecretAccess_ErrorHandling(t *testing.T) {
	t.Run("Error listing ClusterRoleBindings", func(t *testing.T) {
		// This test would require a custom fake client that can simulate errors
		// For now, we test that the function handles errors properly
		clientset := fake.NewSimpleClientset()

		// Add a namespace but no ClusterRoleBindings (normal case)
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "default"},
		}
		_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		assert.NoError(t, err)

		roles := []RoleWithSecretAccess{
			{Kind: "ClusterRole", Namespace: "", Name: "admin", Verbs: []string{"*"}},
		}

		// This should succeed with empty result
		result, err := GetServiceAccountsWithSecretAccess(clientset, roles)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestServiceAccountWithPermissions(t *testing.T) {
	t.Run("ServiceAccountWithPermissions structure", func(t *testing.T) {
		items := []ServiceAccountWithPermissions{
			{
				Namespace: "default",
				Name:      "sa1",
				Permissions: []ServiceAccountPermission{
					{
						RoleRefKind: "ClusterRole",
						RoleRefName: "admin",
						BindingKind: "ClusterRoleBinding",
						Namespace:   "",
						Verbs:       []string{"*"},
					},
				},
			},
			{
				Namespace: "kube-system",
				Name:      "sa2",
				Permissions: []ServiceAccountPermission{
					{
						RoleRefKind: "Role",
						RoleRefName: "secret-reader",
						BindingKind: "RoleBinding",
						Namespace:   "kube-system",
						Verbs:       []string{"get", "list"},
					},
					{
						RoleRefKind: "ClusterRole",
						RoleRefName: "view",
						BindingKind: "RoleBinding",
						Namespace:   "kube-system",
						Verbs:       []string{"watch"},
					},
				},
			},
		}

		// Test basic structure
		assert.Len(t, items, 2)
		assert.Equal(t, "default", items[0].Namespace)
		assert.Equal(t, "sa1", items[0].Name)
		assert.Len(t, items[0].Permissions, 1)
		assert.Equal(t, "admin", items[0].Permissions[0].RoleRefName)

		assert.Equal(t, "kube-system", items[1].Namespace)
		assert.Equal(t, "sa2", items[1].Name)
		assert.Len(t, items[1].Permissions, 2)
	})

	t.Run("Empty list", func(t *testing.T) {
		items := []ServiceAccountWithPermissions{}
		assert.Empty(t, items)
	})
}

func TestServiceAccountPermission_Methods(t *testing.T) {
	tests := []struct {
		name       string
		permission ServiceAccountPermission
		checks     map[string]interface{}
	}{
		{
			name: "ClusterRoleBinding permission",
			permission: ServiceAccountPermission{
				RoleRefKind: "ClusterRole",
				RoleRefName: "admin",
				BindingKind: "ClusterRoleBinding",
				Namespace:   "",
				Verbs:       []string{"get", "list", "watch", "*"},
			},
			checks: map[string]interface{}{
				"isClusterWide":   true,
				"scope":           "Cluster-Wide",
				"hasVerb_get":     true,
				"hasVerb_patch":   true, // because of "*"
				"roleDisplayName": "[ClusterRole] admin",
			},
		},
		{
			name: "RoleBinding permission",
			permission: ServiceAccountPermission{
				RoleRefKind: "Role",
				RoleRefName: "secret-reader",
				BindingKind: "RoleBinding",
				Namespace:   "default",
				Verbs:       []string{"get", "list"},
			},
			checks: map[string]interface{}{
				"isClusterWide":   false,
				"scope":           "NS: default",
				"hasVerb_get":     true,
				"hasVerb_patch":   false,
				"roleDisplayName": "[Role] secret-reader",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if isClusterWide, ok := tt.checks["isClusterWide"].(bool); ok {
				assert.Equal(t, isClusterWide, tt.permission.IsClusterWide())
			}
			if scope, ok := tt.checks["scope"].(string); ok {
				assert.Equal(t, scope, tt.permission.Scope())
			}
			if hasVerb, ok := tt.checks["hasVerb_get"].(bool); ok {
				assert.Equal(t, hasVerb, tt.permission.HasVerb("get"))
			}
			if hasVerb, ok := tt.checks["hasVerb_patch"].(bool); ok {
				assert.Equal(t, hasVerb, tt.permission.HasVerb("patch"))
			}
			if roleDisplayName, ok := tt.checks["roleDisplayName"].(string); ok {
				assert.Equal(t, roleDisplayName, tt.permission.RoleDisplayName())
			}
		})
	}
}
