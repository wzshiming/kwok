/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package image

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	containerregistryv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	// annotationRefName is the annotation read by podman and Apple container to name the image.
	annotationRefName = "org.opencontainers.image.ref.name"
	// annotationContainerdImageName is the annotation read by containerd based runtimes to name the image.
	annotationContainerdImageName = "io.containerd.image.name"

	ociLayoutFile      = "oci-layout"
	ociLayoutContent   = `{"imageLayoutVersion":"1.0.0"}`
	ociIndexFile       = "index.json"
	dockerManifestFile = "manifest.json"
)

// dockerManifest is an entry of the manifest.json of the `docker save` format.
type dockerManifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

// writeArchive writes the image to dest as an OCI image layout archive
// that also carries the manifest.json of the `docker save` format,
// so that it can be loaded by runtimes that only accept one of the two.
func writeArchive(img containerregistryv1.Image, ref name.Reference, dest string) (err error) {
	tmp := dest + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = f.Close()
			_ = os.Remove(tmp)
		}
	}()

	err = writeArchiveTo(f, img, ref)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmp, dest)
}

func writeArchiveTo(w io.Writer, img containerregistryv1.Image, ref name.Reference) error {
	rawManifest, err := img.RawManifest()
	if err != nil {
		return err
	}
	manifestDigest, err := img.Digest()
	if err != nil {
		return err
	}
	mediaType, err := img.MediaType()
	if err != nil {
		return err
	}
	rawConfig, err := img.RawConfigFile()
	if err != nil {
		return err
	}
	configDigest, err := img.ConfigName()
	if err != nil {
		return err
	}
	layers, err := img.Layers()
	if err != nil {
		return err
	}

	refName := normalizeRefName(ref)
	annotations := map[string]string{
		annotationContainerdImageName: refName,
	}
	var repoTags []string
	if _, ok := ref.(name.Tag); ok {
		annotations[annotationRefName] = refName
		repoTags = []string{refName}
	}

	desc := containerregistryv1.Descriptor{
		MediaType:   mediaType,
		Size:        int64(len(rawManifest)),
		Digest:      manifestDigest,
		Annotations: annotations,
	}
	if cf, err := img.ConfigFile(); err == nil && cf != nil && cf.OS != "" {
		desc.Platform = cf.Platform()
	}
	rawIndex, err := json.Marshal(containerregistryv1.IndexManifest{
		SchemaVersion: 2,
		MediaType:     types.OCIImageIndex,
		Manifests:     []containerregistryv1.Descriptor{desc},
	})
	if err != nil {
		return err
	}

	tw := tar.NewWriter(w)

	err = writeTarFile(tw, ociLayoutFile, int64(len(ociLayoutContent)), strings.NewReader(ociLayoutContent))
	if err != nil {
		return err
	}
	err = writeTarFile(tw, ociIndexFile, int64(len(rawIndex)), bytes.NewReader(rawIndex))
	if err != nil {
		return err
	}

	written := map[containerregistryv1.Hash]struct{}{}
	writeBlob := func(digest containerregistryv1.Hash, size int64, r io.Reader) error {
		if _, ok := written[digest]; ok {
			return nil
		}
		written[digest] = struct{}{}
		return writeTarFile(tw, blobPath(digest), size, r)
	}

	err = writeBlob(manifestDigest, int64(len(rawManifest)), bytes.NewReader(rawManifest))
	if err != nil {
		return err
	}
	err = writeBlob(configDigest, int64(len(rawConfig)), bytes.NewReader(rawConfig))
	if err != nil {
		return err
	}

	layerPaths := make([]string, 0, len(layers))
	for _, layer := range layers {
		digest, err := layer.Digest()
		if err != nil {
			return err
		}
		layerPaths = append(layerPaths, blobPath(digest))
		if _, ok := written[digest]; ok {
			continue
		}
		size, err := layer.Size()
		if err != nil {
			return err
		}
		rc, err := layer.Compressed()
		if err != nil {
			return err
		}
		err = writeBlob(digest, size, rc)
		_ = rc.Close()
		if err != nil {
			return err
		}
	}

	rawDockerManifest, err := json.Marshal([]dockerManifest{
		{
			Config:   blobPath(configDigest),
			RepoTags: repoTags,
			Layers:   layerPaths,
		},
	})
	if err != nil {
		return err
	}
	err = writeTarFile(tw, dockerManifestFile, int64(len(rawDockerManifest)), bytes.NewReader(rawDockerManifest))
	if err != nil {
		return err
	}

	return tw.Close()
}

func writeTarFile(tw *tar.Writer, name string, size int64, r io.Reader) error {
	err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     name,
		Size:     size,
		Mode:     0644,
	})
	if err != nil {
		return err
	}
	n, err := io.Copy(tw, r)
	if err != nil {
		return fmt.Errorf("writing %s: %w", name, err)
	}
	if n != size {
		return fmt.Errorf("writing %s: wrote %d bytes, expected %d", name, n, size)
	}
	return nil
}

func blobPath(digest containerregistryv1.Hash) string {
	return path.Join("blobs", digest.Algorithm, digest.Hex)
}

// normalizeRefName returns the fully qualified name of the reference,
// with the registry of Docker Hub spelled the way container runtimes expect it.
func normalizeRefName(ref name.Reference) string {
	n := ref.Name()
	if rest, ok := strings.CutPrefix(n, name.DefaultRegistry+"/"); ok {
		return "docker.io/" + rest
	}
	return n
}
