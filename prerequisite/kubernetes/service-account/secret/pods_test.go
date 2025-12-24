package secret

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	kubernetesSecret "github.com/ctrsploit/ctrsploit/pkg/kubernetes/secret"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
)

const (
	// Timeouts for resource operations
	podReadyTimeout        = 60 * time.Second
	rbacPropagationDelay   = 2 * time.Second
	rbacPropagationMaxWait = 30 * time.Second
	cleanupTimeout         = 30 * time.Second
	pollInterval           = 1 * time.Second
)

// testResources holds all Kubernetes resources created for a test scenario
type testResources struct {
	namespace      string
	serviceAccount string
	pod            string
	role           string
	roleBinding    string
}

func TestE2E_PodsWithSecretAccess(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		createRBAC  bool
		expectFound bool
	}{
		"with_secret_access": {
			createRBAC:  true,
			expectFound: true,
		},
		"without_secret_access": {
			createRBAC:  false,
			expectFound: false,
		},
	}

	clientset, err := kubernetes.GetKubernetesClient()
	require.NoError(t, err, "Failed to get Kubernetes client")

	// If TEST_ENV is set, run only that specific test case
	// Otherwise, run all test cases
	var testCasesToRun map[string]struct {
		createRBAC  bool
		expectFound bool
	}
	if testEnv != "" {
		test, ok := tests[testEnv]
		if !ok {
			t.Skipf("Skipping test for unsupported environment: %s", testEnv)
		}
		testCasesToRun = map[string]struct {
			createRBAC  bool
			expectFound bool
		}{testEnv: test}
	} else {
		testCasesToRun = tests
	}

	for testName, test := range testCasesToRun {
		t.Run(testName, func(t *testing.T) {
			// Generate unique namespace name based on test name
			testNS := "e2e-test-secret-pods-" + strings.ToLower(strings.ReplaceAll(testName, "_", "-"))
			resources := &testResources{
				namespace:      testNS,
				serviceAccount: "test-sa",
				pod:            "test-pod",
				role:           "secret-reader",
				roleBinding:    "secret-reader-binding",
			}

			// Setup test resources
			ctx := context.Background()
			setupTestResources(t, ctx, clientset, resources, test.createRBAC)

			// Cleanup with proper error handling
			defer func() {
				cleanupCtx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
				defer cancel()
				cleanupTestResources(t, cleanupCtx, clientset, resources)
			}()

			// Wait for RBAC to propagate if needed
			if test.createRBAC {
				waitForRBACPropagation(t, ctx, clientset, resources)
			}

			// Run the prerequisite check
			satisfied, err := HasPodsWithSecretAccess.Check()
			require.NoError(t, err, "Failed to check prerequisite")

			// Extract pods
			pods, err := GetPodsWithSecretAccess(&HasPodsWithSecretAccess)
			require.NoError(t, err, "Failed to get pods with secret access")

			// Print all pods with secret access found
			printPodsWithSecretAccess(t, pods)

			// Verify expectations
			verifyTestResults(t, resources, pods, satisfied, test.expectFound)
		})
	}
}

// setupTestResources creates all necessary Kubernetes resources for a test scenario
func setupTestResources(t *testing.T, ctx context.Context, clientset *k8sclient.Clientset, resources *testResources, createRBAC bool) {
	// Create namespace
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: resources.namespace,
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create namespace %s: %v", resources.namespace, err)
	}
	t.Logf("Created namespace: %s", resources.namespace)

	// Create ServiceAccount
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resources.serviceAccount,
			Namespace: resources.namespace,
		},
	}
	_, err = clientset.CoreV1().ServiceAccounts(resources.namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create ServiceAccount %s/%s: %v", resources.namespace, resources.serviceAccount, err)
	}
	t.Logf("Created ServiceAccount: %s/%s", resources.namespace, resources.serviceAccount)

	// Create RBAC resources if needed
	if createRBAC {
		createRBACResources(t, ctx, clientset, resources)
	}

	// Create Pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resources.pod,
			Namespace: resources.namespace,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: resources.serviceAccount,
			Containers: []corev1.Container{
				{
					Name:    "test-container",
					Image:   "busybox:latest",
					Command: []string{"sleep", "3600"},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	_, err = clientset.CoreV1().Pods(resources.namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create Pod %s/%s: %v", resources.namespace, resources.pod, err)
	}
	t.Logf("Created Pod: %s/%s", resources.namespace, resources.pod)

	// Wait for pod to be ready
	podCtx, cancel := context.WithTimeout(ctx, podReadyTimeout)
	defer cancel()
	err = waitForPodReady(podCtx, clientset, resources.namespace, resources.pod)
	if err != nil {
		t.Logf("Warning: Pod %s/%s not ready within timeout: %v", resources.namespace, resources.pod, err)
	} else {
		t.Logf("Pod %s/%s is ready", resources.namespace, resources.pod)
	}
}

// createRBACResources creates Role and RoleBinding for secret access
func createRBACResources(t *testing.T, ctx context.Context, clientset *k8sclient.Clientset, resources *testResources) {
	// Create Role with secret access
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resources.role,
			Namespace: resources.namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
	_, err := clientset.RbacV1().Roles(resources.namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create Role %s/%s: %v", resources.namespace, resources.role, err)
	}
	t.Logf("Created Role: %s/%s", resources.namespace, resources.role)

	// Create RoleBinding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resources.roleBinding,
			Namespace: resources.namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     resources.role,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      resources.serviceAccount,
				Namespace: resources.namespace,
			},
		},
	}
	_, err = clientset.RbacV1().RoleBindings(resources.namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create RoleBinding %s/%s: %v", resources.namespace, resources.roleBinding, err)
	}
	t.Logf("Created RoleBinding: %s/%s", resources.namespace, resources.roleBinding)
}

// cleanupTestResources deletes all test resources in reverse order of creation
func cleanupTestResources(t *testing.T, ctx context.Context, clientset *k8sclient.Clientset, resources *testResources) {
	t.Logf("Cleaning up test resources in namespace: %s", resources.namespace)

	// Delete Pod first
	err := clientset.CoreV1().Pods(resources.namespace).Delete(ctx, resources.pod, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		t.Logf("Warning: Failed to delete Pod %s/%s: %v", resources.namespace, resources.pod, err)
	} else if err == nil {
		// Wait for pod deletion
		waitForResourceDeletion(ctx, func() error {
			_, err := clientset.CoreV1().Pods(resources.namespace).Get(ctx, resources.pod, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		})
	}

	// Delete RoleBinding
	err = clientset.RbacV1().RoleBindings(resources.namespace).Delete(ctx, resources.roleBinding, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		t.Logf("Warning: Failed to delete RoleBinding %s/%s: %v", resources.namespace, resources.roleBinding, err)
	}

	// Delete Role
	err = clientset.RbacV1().Roles(resources.namespace).Delete(ctx, resources.role, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		t.Logf("Warning: Failed to delete Role %s/%s: %v", resources.namespace, resources.role, err)
	}

	// Delete ServiceAccount
	err = clientset.CoreV1().ServiceAccounts(resources.namespace).Delete(ctx, resources.serviceAccount, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		t.Logf("Warning: Failed to delete ServiceAccount %s/%s: %v", resources.namespace, resources.serviceAccount, err)
	}

	// Delete namespace (this will cascade delete remaining resources)
	err = clientset.CoreV1().Namespaces().Delete(ctx, resources.namespace, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		t.Logf("Warning: Failed to delete namespace %s: %v", resources.namespace, err)
	} else if err == nil {
		// Wait for namespace deletion
		waitForResourceDeletion(ctx, func() error {
			_, err := clientset.CoreV1().Namespaces().Get(ctx, resources.namespace, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		})
		t.Logf("Successfully cleaned up namespace: %s", resources.namespace)
	}
}

// waitForResourceDeletion waits for a resource to be deleted
func waitForResourceDeletion(ctx context.Context, checkFn func() error) {
	deadline := time.Now().Add(cleanupTimeout)
	for time.Now().Before(deadline) {
		if err := checkFn(); err == nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(pollInterval):
			// Continue polling
		}
	}
}

// waitForRBACPropagation waits for RBAC permissions to propagate
// RBAC propagation in Kubernetes can take a few seconds, so we wait and verify
// that the RoleBinding exists before proceeding with the test
func waitForRBACPropagation(t *testing.T, ctx context.Context, clientset *k8sclient.Clientset, resources *testResources) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	deadline := time.Now().Add(rbacPropagationMaxWait)
	t.Logf("Waiting for RBAC to propagate for ServiceAccount %s/%s...", resources.namespace, resources.serviceAccount)

	// Initial delay to allow RBAC to start propagating
	select {
	case <-ctx.Done():
		return
	case <-time.After(rbacPropagationDelay):
	}

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Verify that RoleBinding exists and is accessible
			_, err := clientset.RbacV1().RoleBindings(resources.namespace).Get(ctx, resources.roleBinding, metav1.GetOptions{})
			if err == nil {
				// Also verify Role exists
				_, err = clientset.RbacV1().Roles(resources.namespace).Get(ctx, resources.role, metav1.GetOptions{})
				if err == nil {
					t.Logf("RBAC resources are ready and accessible")
					return
				}
			}
		}
	}
	t.Logf("Warning: RBAC propagation wait completed (may not be fully propagated)")
}

// verifyTestResults verifies the test results against expectations
func verifyTestResults(t *testing.T, resources *testResources, pods []kubernetesSecret.PodWithSecretAccess, satisfied bool, expectFound bool) {
	// Check if our test pod is found
	found := findPodInList(pods, resources.namespace, resources.pod, resources.serviceAccount)

	if expectFound {
		// Should find the pod with secret access
		require.True(t, satisfied, "Should find pods with secret access")
		require.GreaterOrEqual(t, len(pods), 1, "Expected at least 1 pod with secret access")

		if found {
			t.Logf("Found expected pod: %s/%s with SA: %s", resources.namespace, resources.pod, resources.serviceAccount)
		} else {
			logAllPods(t, pods)
			t.Errorf("Expected pod %s/%s not found in results", resources.namespace, resources.pod)
		}
	} else {
		// Should NOT find the pod (it has no secret access)
		if found {
			t.Errorf("Unexpectedly found pod %s/%s with secret access (should not have access)", resources.namespace, resources.pod)
		} else {
			t.Logf("Test pod %s/%s correctly not found in pods with secret access", resources.namespace, resources.pod)
		}
	}
}

// findPodInList checks if a specific pod exists in the list
func findPodInList(pods []kubernetesSecret.PodWithSecretAccess, namespace, podName, saName string) bool {
	for _, p := range pods {
		if p.Name == podName && p.Namespace == namespace && p.ServiceAccountName == saName {
			return true
		}
	}
	return false
}

// printPodsWithSecretAccess prints all pods with secret access in a detailed format
func printPodsWithSecretAccess(t *testing.T, pods []kubernetesSecret.PodWithSecretAccess) {
	if len(pods) == 0 {
		t.Logf("No pods with secret access found")
		return
	}

	t.Logf("=== Found %d pod(s) with secret access ===", len(pods))
	for i, p := range pods {
		t.Logf("[%d] Pod: %s/%s", i+1, p.Namespace, p.Name)
		t.Logf("     ServiceAccount: %s", p.ServiceAccountName)
		if len(p.Permissions) > 0 {
			t.Logf("     Permissions:")
			for j, perm := range p.Permissions {
				bindingInfo := perm.BindingKind
				if perm.Namespace != "" {
					bindingInfo += " (" + perm.Namespace + ")"
				}
				t.Logf("       [%d] %s: %s via %s (Verbs: %v)", j+1, perm.RoleRefKind, perm.RoleRefName, bindingInfo, perm.Verbs)
			}
		}
	}
	t.Logf("==========================================")
}

// logAllPods logs all pods found for debugging (simplified format)
func logAllPods(t *testing.T, pods []kubernetesSecret.PodWithSecretAccess) {
	var podNames []string
	for _, p := range pods {
		podNames = append(podNames, p.Namespace+"/"+p.Name+" (SA: "+p.ServiceAccountName+")")
	}
	t.Logf("Found pods with secret access: %s", strings.Join(podNames, ", "))
}

// waitForPodReady waits for a pod to be in Ready condition or Running phase
func waitForPodReady(ctx context.Context, clientset *k8sclient.Clientset, namespace, podName string) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for pod %s/%s: %w", namespace, podName, ctx.Err())
		case <-ticker.C:
			pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return fmt.Errorf("pod %s/%s not found", namespace, podName)
				}
				continue
			}

			// Check if pod is in a terminal error state
			if pod.Status.Phase == corev1.PodFailed {
				return fmt.Errorf("pod %s/%s failed: %s", namespace, podName, pod.Status.Reason)
			}

			// Check Ready condition
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
					return nil
				}
			}

			// Also accept Running phase (for pods without readiness probe)
			if pod.Status.Phase == corev1.PodRunning {
				return nil
			}
		}
	}
}
