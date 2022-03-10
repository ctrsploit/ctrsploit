module github.com/ctrsploit/ctrsploit

go 1.16

require (
	github.com/containerd/cgroups v1.0.1
	github.com/containerd/containerd v1.5.2
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v20.10.7+incompatible
	github.com/fatih/color v1.7.0
	github.com/google/cadvisor v0.40.0
	github.com/google/uuid v1.2.0
	github.com/mitchellh/gox v1.0.1 // indirect
	github.com/moby/sys/mountinfo v0.4.1
	github.com/opencontainers/runc v1.0.0-rc95
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/ssst0n3/awesome_libs v0.6.7
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
	k8s.io/client-go v0.20.6
)
