package image

import (
	"context"
	"fmt"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ResolveDiffIDs fetches the OCI image config for ref and returns the
// layer diff_ids from config.rootfs.diff_ids.
//
// ref is a standard image reference such as "alpine:latest",
// "ubuntu:22.04", "docker.io/library/nginx:alpine".
func ResolveDiffIDs(ref string) ([]v1.Hash, error) {
	parsed, err := name.ParseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("parse image reference %q: %w", ref, err)
	}

	img, err := remote.Image(parsed, remote.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("fetch image %q: %w", ref, err)
	}

	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("read image config for %q: %w", ref, err)
	}

	return cfg.RootFS.DiffIDs, nil
}
