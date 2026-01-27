package oci

import (
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	RegistryPath = "/v2"
	// ComposeProjectArtifactType is the OCI 1.1-compliant artifact type value
	// for the generated image manifest.
	// https://github.com/docker/compose/blob/v2.40.1/internal/oci/push.go#L41C2-L41C71
	ComposeProjectArtifactType = "application/vnd.docker.compose.project"
	// ComposeYAMLMediaType is the media type for each layer (Compose file)
	// in the image manifest.
	// https://github.com/docker/compose/blob/v2.40.1/internal/oci/push.go#L44
	ComposeYAMLMediaType = "application/vnd.docker.compose.file+yaml"
	// ComposeEmptyConfigMediaType is a media type used for the config descriptor
	// https://github.com/docker/compose/blob/v2.40.1/internal/oci/push.go#L57
	ComposeEmptyConfigMediaType = "application/vnd.docker.compose.config.empty.v1+json"
)

// DescriptorForComposeFile creates the main compose file layer descriptor
// Reference: https://github.com/docker/compose/blob/v2.40.1/internal/oci/push.go#L70
func DescriptorForComposeFile(path string, content []byte) v1.Descriptor {
	return v1.Descriptor{
		MediaType: ComposeYAMLMediaType,
		Digest:    digest.FromBytes(content),
		Size:      int64(len(content)),
		Annotations: map[string]string{
			"com.docker.compose.version": "v2",
			"com.docker.compose.file":    path,
		},
		// Data: content,
	}
}
