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

func TestGetPodsWithSecretAccess(t *testing.T) {
	tests := []struct {
		name           string
		saWithPerms    []ServiceAccountWithPermissions
		namespaces     []*corev1.Namespace
		podsByNS       map[string][]*corev1.Pod
		expectedResult []PodWithSecretAccess
		expectError    bool
	}{
		{
			name: "Pod with ServiceAccount that has secret access",
			saWithPerms: []ServiceAccountWithPermissions{
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "sa1",
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
			name: "Pod with ServiceAccount that has no secret access",
			saWithPerms: []ServiceAccountWithPermissions{
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa2", // Different SA without permissions
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{},
			expectError:    false,
		},
		{
			name: "Pod without ServiceAccount specified (defaults to 'default')",
			saWithPerms: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "default", // Default ServiceAccount
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							// ServiceAccountName is empty, should default to "default"
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "default",
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
			name: "Pod with ServiceAccount that has multiple permissions",
			saWithPerms: []ServiceAccountWithPermissions{
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "sa1",
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
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "sa1",
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
			name: "Multiple Pods using same ServiceAccount",
			saWithPerms: []ServiceAccountWithPermissions{
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod2",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "sa1",
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
					Namespace:          "default",
					Name:               "pod2",
					ServiceAccountName: "sa1",
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
			name: "Pods in multiple namespaces",
			saWithPerms: []ServiceAccountWithPermissions{
				{
					Namespace: "default",
					Name:      "sa1",
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
				"kube-system": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod2",
							Namespace: "kube-system",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa2",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "sa1",
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
				{
					Namespace:          "kube-system",
					Name:               "pod2",
					ServiceAccountName: "sa2",
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
			name:        "Empty ServiceAccount permissions list",
			saWithPerms: []ServiceAccountWithPermissions{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{},
			expectError:    false,
		},
		{
			name: "ServiceAccount with empty permissions",
			saWithPerms: []ServiceAccountWithPermissions{
				{
					Namespace:   "default",
					Name:        "sa1",
					Permissions: []ServiceAccountPermission{}, // Empty permissions
				},
			},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{},
			expectError:    false,
		},
		{
			name: "No Pods in namespace",
			saWithPerms: []ServiceAccountWithPermissions{
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS:       map[string][]*corev1.Pod{},
			expectedResult: []PodWithSecretAccess{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake clientset
			clientset := fake.NewSimpleClientset()

			// Add namespaces
			for _, ns := range tt.namespaces {
				_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add Pods by namespace
			for ns, pods := range tt.podsByNS {
				for _, pod := range pods {
					_, err := clientset.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
					assert.NoError(t, err)
				}
			}

			// Test the function
			result, err := GetPodsWithSecretAccess(clientset, tt.saWithPerms)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// Compare results - use ElementsMatch since order doesn't matter
				assert.Equal(t, len(tt.expectedResult), len(result), "Result count mismatch")
				assert.ElementsMatch(t, tt.expectedResult, result, "Pods with secret access mismatch")
			}
		})
	}
}

func TestGetPodsWithSecretAccess_ErrorHandling(t *testing.T) {
	t.Run("Empty clientset with no namespaces", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()

		saWithPerms := []ServiceAccountWithPermissions{
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
		}

		// This should succeed with empty result (no namespaces)
		result, err := GetPodsWithSecretAccess(clientset, saWithPerms)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestPodWithSecretAccess(t *testing.T) {
	t.Run("PodWithSecretAccess structure", func(t *testing.T) {
		items := []PodWithSecretAccess{
			{
				Namespace:          "default",
				Name:               "pod1",
				ServiceAccountName: "sa1",
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
				Namespace:          "kube-system",
				Name:               "pod2",
				ServiceAccountName: "sa2",
				Permissions: []ServiceAccountPermission{
					{
						RoleRefKind: "Role",
						RoleRefName: "secret-reader",
						BindingKind: "RoleBinding",
						Namespace:   "kube-system",
						Verbs:       []string{"get", "list"},
					},
				},
			},
		}

		// Test basic structure
		assert.Len(t, items, 2)
		assert.Equal(t, "default", items[0].Namespace)
		assert.Equal(t, "pod1", items[0].Name)
		assert.Equal(t, "sa1", items[0].ServiceAccountName)
		assert.Len(t, items[0].Permissions, 1)
		assert.Equal(t, "admin", items[0].Permissions[0].RoleRefName)

		assert.Equal(t, "kube-system", items[1].Namespace)
		assert.Equal(t, "pod2", items[1].Name)
		assert.Equal(t, "sa2", items[1].ServiceAccountName)
		assert.Len(t, items[1].Permissions, 1)
		assert.Equal(t, "secret-reader", items[1].Permissions[0].RoleRefName)
	})

	t.Run("Empty list", func(t *testing.T) {
		items := []PodWithSecretAccess{}
		assert.Empty(t, items)
	})
}

func TestGetAllPodsWithSecretAccess(t *testing.T) {
	tests := []struct {
		name                string
		clusterRoles        []*rbacv1.ClusterRole
		rolesByNS           map[string][]*rbacv1.Role
		clusterRoleBindings []*rbacv1.ClusterRoleBinding
		roleBindingsByNS    map[string][]*rbacv1.RoleBinding
		namespaces          []*corev1.Namespace
		podsByNS            map[string][]*corev1.Pod
		expectedResult      []PodWithSecretAccess
		expectError         bool
	}{
		{
			name: "Complete chain: ClusterRole -> ClusterRoleBinding -> ServiceAccount -> Pod",
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
			rolesByNS: map[string][]*rbacv1.Role{},
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
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa1",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod1",
					ServiceAccountName: "sa1",
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
			name: "Complete chain: Role -> RoleBinding -> ServiceAccount -> Pod",
			clusterRoles:        []*rbacv1.ClusterRole{},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{},
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
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod2",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa2",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod2",
					ServiceAccountName: "sa2",
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
			name: "Complete chain with multiple Pods",
			clusterRoles: []*rbacv1.ClusterRole{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "view"},
					Rules: []rbacv1.PolicyRule{
						{
							Resources: []string{"secrets"},
							Verbs:     []string{"get", "list", "watch"},
						},
					},
				},
			},
			rolesByNS: map[string][]*rbacv1.Role{},
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
							Namespace: "default",
							Name:      "sa3",
						},
					},
				},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod3",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa3",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod4",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa3",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{
				{
					Namespace:          "default",
					Name:               "pod3",
					ServiceAccountName: "sa3",
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
					Namespace:          "default",
					Name:               "pod4",
					ServiceAccountName: "sa3",
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
			},
			expectError: false,
		},
		{
			name: "No Pods with secret access (Role without secret access)",
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
			rolesByNS: map[string][]*rbacv1.Role{},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{
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
							Name:      "sa4",
						},
					},
				},
			},
			roleBindingsByNS: map[string][]*rbacv1.RoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS: map[string][]*corev1.Pod{
				"default": {
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pod5",
							Namespace: "default",
						},
						Spec: corev1.PodSpec{
							ServiceAccountName: "sa4",
						},
					},
				},
			},
			expectedResult: []PodWithSecretAccess{},
			expectError:    false,
		},
		{
			name:                "Empty cluster (no roles, no pods)",
			clusterRoles:        []*rbacv1.ClusterRole{},
			rolesByNS:           map[string][]*rbacv1.Role{},
			clusterRoleBindings: []*rbacv1.ClusterRoleBinding{},
			roleBindingsByNS:    map[string][]*rbacv1.RoleBinding{},
			namespaces: []*corev1.Namespace{
				{ObjectMeta: metav1.ObjectMeta{Name: "default"}},
			},
			podsByNS:       map[string][]*corev1.Pod{},
			expectedResult: []PodWithSecretAccess{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake clientset
			clientset := fake.NewSimpleClientset()

			// Add namespaces
			for _, ns := range tt.namespaces {
				_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add ClusterRoles
			for _, cr := range tt.clusterRoles {
				_, err := clientset.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add Roles by namespace
			for ns, roles := range tt.rolesByNS {
				for _, role := range roles {
					_, err := clientset.RbacV1().Roles(ns).Create(context.TODO(), role, metav1.CreateOptions{})
					assert.NoError(t, err)
				}
			}

			// Add ClusterRoleBindings
			for _, crb := range tt.clusterRoleBindings {
				_, err := clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{})
				assert.NoError(t, err)
			}

			// Add RoleBindings by namespace
			for ns, rbs := range tt.roleBindingsByNS {
				for _, rb := range rbs {
					_, err := clientset.RbacV1().RoleBindings(ns).Create(context.TODO(), rb, metav1.CreateOptions{})
					assert.NoError(t, err)
				}
			}

			// Add Pods by namespace
			for ns, pods := range tt.podsByNS {
				for _, pod := range pods {
					_, err := clientset.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
					assert.NoError(t, err)
				}
			}

			// Test the function
			result, err := GetAllPodsWithSecretAccess(clientset)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// Compare results - use ElementsMatch since order doesn't matter
				assert.Equal(t, len(tt.expectedResult), len(result), "Result count mismatch")
				assert.ElementsMatch(t, tt.expectedResult, result, "Pods with secret access mismatch")
			}
		})
	}
}

func TestGetAllPodsWithSecretAccess_ErrorHandling(t *testing.T) {
	t.Run("Empty clientset", func(t *testing.T) {
		clientset := fake.NewSimpleClientset()

		// This should succeed with empty result
		result, err := GetAllPodsWithSecretAccess(clientset)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}
