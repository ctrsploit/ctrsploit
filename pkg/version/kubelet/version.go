package kubelet

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Versions tries to get the kubelet version by multiple methods.
// It first attempts to get the versions by querying the Kubernetes API.
// If that fails, it falls back to getting the version via the kubelet CLI.
// If both methods fail, it aggregates the errors and returns them.
func Versions() ([]*semver.Version, error) {
	type getter func() ([]*semver.Version, error)
	methods := []getter{
		VersionsByK8sApi,
		versionsByCli,
	}
	var errs []error
	for _, fn := range methods {
		if versions, err := fn(); err == nil {
			return versions, nil
		} else {
			awesome_error.CheckDebug(err)
			errs = append(errs, err)
		}
	}
	return nil, errors.Join(errs...)
}

func VersionsByK8sApi() ([]*semver.Version, error) {
	c, err := kubernetes.GetKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes client: %w", err)
	}
	nodes, err := c.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}
	log.Logger.Debugf("Found %d nodes", len(nodes.Items))
	if len(nodes.Items) == 0 {
		return nil, fmt.Errorf("no nodes found in the cluster")
	}
	var versions []*semver.Version
	for _, node := range nodes.Items {
		kubeletVersion := node.Status.NodeInfo.KubeletVersion
		log.Logger.Debugf("Node %s kubelet version: %s", node.Name, kubeletVersion)
		ver, err := semver.NewVersion(kubeletVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to convert kubelet version to semver for node %s: %w", node.Name, err)
		}
		versions = append(versions, ver)
	}
	return versions, nil
}

// versionsByCli is a helper function makes the code more simple.
func versionsByCli() ([]*semver.Version, error) {
	version, err := VersionByCli()
	if err != nil {
		return nil, err
	}
	return []*semver.Version{version}, nil
}

func VersionByCli() (*semver.Version, error) {
	path, err := exec.LookPath("kubelet")
	if err != nil {
		return nil, fmt.Errorf("failed to find kubelet binary: %w", err)
	}
	output, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute kubelet --version: %w", err)
	}
	var versionStr string
	_, err = fmt.Sscanf(string(output), "Kubernetes v%s", &versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubelet version from output: %w", err)
	}
	versionStr = "v" + versionStr
	ver, err := semver.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert kubelet version to semver: %w", err)
	}
	return ver, nil
}
