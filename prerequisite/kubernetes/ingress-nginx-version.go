package kubernetes

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressNginxVersionConstraint struct {
	prerequisite.BasePrerequisite
	Constraint      string
	matchedPods     []string
	matchedServices []IngressNginxService
}

type PortMapping struct {
	ServicePort    int32
	NodePort       int32
	TargetPort     int32
	TargetPortName string
	Name           string
	Protocol       string
}

type IngressNginxService struct {
	Namespace   string
	Name        string
	ClusterIP   string
	ExternalIP  string
	Ports       []PortMapping
	Version     *semver.Version
	Labels      map[string]string
	Selectors   map[string]string
	ServiceType string
}

func (p *IngressNginxVersionConstraint) Check() (bool, error) {
	return p.CheckTemplate(func() {
		p.Satisfied = true
		cons, err := semver.NewConstraint(p.Constraint)
		if err != nil {
			err = fmt.Errorf("failed to parse constraint %s: %w", p.Constraint, err)
			awesome_error.CheckFatal(err)
		}

		services, err := GetIngressNginxServices()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting ingress-nginx services %w", err))
			return
		}
		if len(services) == 0 {
			p.Err = p.WrapErr(fmt.Errorf("no ingress-nginx controller found"))
			return
		}

		p.matchedPods = nil
		p.matchedServices = nil

		for _, service := range services {
			serviceName := fmt.Sprintf("%s/%s", service.Namespace, service.Name)

			var version *semver.Version
			if versionLabel, exists := service.Labels["app.kubernetes.io/version"]; exists {
				version, err = semver.NewVersion(versionLabel)
				if err != nil {
					log.Logger.Debugf("Failed to parse version from service labels: %v", err)
				}
			}

			if version == nil {
				version = service.Version
			}

			if version == nil {
				log.Logger.Warnf("No version found for service %s", serviceName)
				continue
			}

			if cons.Check(version) {
				p.matchedPods = append(p.matchedPods, serviceName)
				p.matchedServices = append(p.matchedServices, service)
			} else {
				p.Satisfied = false
			}
		}
		return
	})
}

func GetIngressNginxServices() ([]IngressNginxService, error) {
	clientset, err := kubernetes.GetKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes client: %w", err)
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var ingressNginxServices []IngressNginxService

	for _, service := range services.Items {
		labels := service.Labels
		if labels["app.kubernetes.io/name"] != "ingress-nginx" {
			continue
		}

		log.Logger.Debugf("Found ingress-nginx service: %s/%s", service.Namespace, service.Name)

		clusterIP := service.Spec.ClusterIP
		externalIP := ""
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			externalIP = service.Status.LoadBalancer.Ingress[0].IP
		}

		podInfo, err := getPodInfo(service.Namespace, service.Spec.Selector)
		if err != nil {
			log.Logger.Debugf("Failed to get pod info for service %s/%s: %v", service.Namespace, service.Name, err)
		}

		var ports []PortMapping
		for _, port := range service.Spec.Ports {
			var targetPort int32
			var targetPortName string

			if port.TargetPort.Type == 0 {
				targetPort = port.TargetPort.IntVal
			} else {
				targetPortName = port.TargetPort.StrVal
				if podInfo != nil {
					if resolvedPort, ok := podInfo.NamedPorts[targetPortName]; ok {
						targetPort = resolvedPort
					} else {
						log.Logger.Debugf("Named port %s not found in pod info", targetPortName)
						targetPort = 0
					}
				} else {
					log.Logger.Debugf("No pod info available to resolve named port %s", targetPortName)
					targetPort = 0
				}
			}

			portMapping := PortMapping{
				ServicePort:    port.Port,
				NodePort:       port.NodePort,
				TargetPort:     targetPort,
				TargetPortName: targetPortName,
				Name:           port.Name,
				Protocol:       string(port.Protocol),
			}
			ports = append(ports, portMapping)
		}

		var version *semver.Version
		if podInfo != nil && podInfo.Image != "" {
			version, err = extractVersionFromImage(podInfo.Image)
			if err != nil {
				log.Logger.Debugf("Failed to extract version from image %s: %v", podInfo.Image, err)
			} else {
				log.Logger.Debugf("Extracted version %s from image %s", version.String(), podInfo.Image)
			}
		}

		ingressNginxService := IngressNginxService{
			Namespace:   service.Namespace,
			Name:        service.Name,
			ClusterIP:   clusterIP,
			ExternalIP:  externalIP,
			Ports:       ports,
			Version:     version,
			Labels:      service.Labels,
			Selectors:   service.Spec.Selector,
			ServiceType: string(service.Spec.Type),
		}

		ingressNginxServices = append(ingressNginxServices, ingressNginxService)

		portInfo := ""
		for i, port := range ports {
			if i > 0 {
				portInfo += ", "
			}
			portInfo += fmt.Sprintf("%d:%d->%d", port.ServicePort, port.NodePort, port.TargetPort)
		}

		log.Logger.Debugf("Found ingress-nginx controller %s/%s (ClusterIP: %s, ExternalIP: %s, Ports: [%s], Type: %s)",
			service.Namespace, service.Name, clusterIP, externalIP, portInfo, service.Spec.Type)
	}

	return ingressNginxServices, nil
}

type PodInfo struct {
	Image      string
	NamedPorts map[string]int32
}

func getPodInfo(namespace string, selector map[string]string) (*PodInfo, error) {
	clientset, err := kubernetes.GetKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes client: %w", err)
	}

	if len(selector) == 0 {
		return nil, fmt.Errorf("no selector provided")
	}

	var selectorParts []string
	for key, value := range selector {
		selectorParts = append(selectorParts, fmt.Sprintf("%s=%s", key, value))
	}
	labelSelector := strings.Join(selectorParts, ",")

	// Get pods
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods found with selector %s", labelSelector)
	}

	pod := pods.Items[0]
	podInfo := &PodInfo{
		NamedPorts: make(map[string]int32),
	}

	for _, container := range pod.Spec.Containers {
		if strings.Contains(container.Image, "ingress-nginx") || strings.Contains(container.Image, "nginx-ingress") {
			podInfo.Image = container.Image
		}

		for _, port := range container.Ports {
			if port.Name != "" {
				podInfo.NamedPorts[port.Name] = port.ContainerPort
			}
		}
	}

	return podInfo, nil
}

func extractVersionFromImage(image string) (*semver.Version, error) {
	re := regexp.MustCompile(`:v?([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindStringSubmatch(image)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no version found in image: %s", image)
	}

	versionStr := matches[1]
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %s: %w", versionStr, err)
	}

	return version, nil
}

func GetIngressNginxPods(s prerequisite.Set) ([]string, error) {
	_, _ = s.Check()
	if p, ok := s.(*IngressNginxVersionConstraint); ok {
		return p.matchedPods, nil
	}
	return nil, fmt.Errorf("unsupported prerequisite type: %T", s)
}

func GetVulnerableIngressNginxServices(s prerequisite.Set) ([]IngressNginxService, error) {
	_, _ = s.Check()
	if p, ok := s.(*IngressNginxVersionConstraint); ok {
		return p.matchedServices, nil
	}
	return nil, fmt.Errorf("unsupported prerequisite type: %T", s)
}

var (
	ConstraintIngressNginxVulnerableToCVE_2021_25748 = "< 1.2.1"
	VulnerableToCVE_2021_25748                       = IngressNginxVersionConstraint{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "cve-2021-25748-vulnerable",
			Info:   "ingress-nginx < 1.2.1 (vulnerable to CVE-2021-25748)",
			ExeEnv: exeenv.InHost | exeenv.InContainer,
		},
		Constraint: ConstraintIngressNginxVulnerableToCVE_2021_25748,
	}
	MaybeVulnerableToCVE_2021_25748 = &VulnerableToCVE_2021_25748
)
