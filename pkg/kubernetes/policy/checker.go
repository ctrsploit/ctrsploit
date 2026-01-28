package policy

import (
	"context"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckDangerousPermissions detects high-risk permissions of the current token
func CheckDangerousPermissions(clientset kubernetes.Interface, namespace string, permissions []DangerousPermission) ([]CheckResult, error) {
	var results []CheckResult

	for _, perm := range permissions {
		for _, verb := range perm.Verbs {
			allowed, err := CanI(clientset, verb, perm.Group, perm.Resource, perm.Subresource, namespace)
			if err != nil {
				continue
			}
			if allowed {
				results = append(results, CheckResult{
					Permission:  perm,
					Allowed:     true,
					Namespace:   namespace,
					MatchedVerb: verb,
				})
				break // Only record the first matched verb
			}
		}
	}

	return results, nil
}

// CanI uses SelfSubjectAccessReview to detect permissions
func CanI(clientset kubernetes.Interface, verb, group, resource, subresource, namespace string) (bool, error) {
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace:   namespace,
				Verb:        verb,
				Group:       group,
				Resource:    resource,
				Subresource: subresource,
			},
		},
	}

	response, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(
		context.TODO(), sar, metav1.CreateOptions{},
	)
	if err != nil {
		return false, err
	}

	return response.Status.Allowed, nil
}
