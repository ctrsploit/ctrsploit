package image

import (
	"encoding/json"
	"github.com/containerd/containerd/images"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"runtime"
	"time"
)

func Digest(b []byte) (digest.Digest, int64) {
	return digest.FromBytes(b), int64(len(b))
}

type Image struct {
	tar   *Tar
	layer []byte
}

func NewImage(layer []byte) *Image {
	return &Image{
		layer: layer,
		tar:   NewTar(),
	}
}

func (i *Image) Build(imageName string) (bytes []byte, err error) {
	layout, err := json.Marshal(i.layout())
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	layerDigest, layerSize := Digest(i.layer)
	config, err := json.Marshal(i.config(layerDigest))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	configDigest, configSize := Digest(config)
	manifest, err := json.Marshal(i.manifest(configDigest, configSize, layerDigest, layerSize))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	manifestDigest, manifestSize := Digest(manifest)
	index, err := json.Marshal(i.index(imageName, manifestDigest, manifestSize))

	if err = i.tar.Dir("blobs/"); err != nil {
		return
	}
	if err = i.tar.Dir("blobs/sha256/"); err != nil {
		return
	}
	if err = i.tar.Reg("oci-layout", layout); err != nil {
		return
	}
	if err = i.tar.Reg("index.json", index); err != nil {
		return
	}
	if err = i.tar.Reg("blobs/sha256/"+configDigest.Hex(), config); err != nil {
		return
	}
	if err = i.tar.Reg("blobs/sha256/"+manifestDigest.Hex(), manifest); err != nil {
		return
	}
	if err = i.tar.Reg("blobs/sha256/"+layerDigest.Hex(), i.layer); err != nil {
		return
	}
	bytes, err = i.tar.Build()
	if err != nil {
		return
	}
	return
}

func (i *Image) layout() (layout v1.ImageLayout) {
	return v1.ImageLayout{
		Version: v1.ImageLayoutVersion,
	}
}

func (i *Image) index(imageName string, manifestDigest digest.Digest, manifestSize int64) (index v1.Index) {
	return v1.Index{
		Versioned: specs.Versioned{
			SchemaVersion: 2,
		},
		MediaType: v1.MediaTypeImageIndex,
		Manifests: []v1.Descriptor{
			{
				MediaType: v1.MediaTypeImageManifest,
				Digest:    manifestDigest,
				Size:      manifestSize,
				Annotations: map[string]string{
					images.AnnotationImageName: imageName,
					v1.AnnotationRefName:       "latest",
				},
			},
		},
	}
}

func (i *Image) manifest(configDigest digest.Digest, configSize int64, layerDigest digest.Digest, layerSize int64) (manifest v1.Manifest) {
	return v1.Manifest{
		Versioned: specs.Versioned{
			SchemaVersion: 2,
		},
		MediaType: v1.MediaTypeImageManifest,
		Config: v1.Descriptor{
			MediaType: v1.MediaTypeImageConfig,
			Digest:    configDigest,
			Size:      configSize,
		},
		Layers: []v1.Descriptor{
			{
				MediaType: v1.MediaTypeImageLayer,
				Digest:    layerDigest,
				Size:      layerSize,
			},
		},
	}
}

func (i *Image) config(layerDigest digest.Digest) (err v1.Image) {
	now := time.Now()
	return v1.Image{
		Created: &now,
		Platform: v1.Platform{
			Architecture: runtime.GOARCH,
			OS:           runtime.GOOS,
		},
		RootFS: v1.RootFS{
			Type:    "layers",
			DiffIDs: []digest.Digest{layerDigest},
		},
	}
}
