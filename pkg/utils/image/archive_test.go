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
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	containerregistryv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/validate"
)

// untar extracts the archive to a directory and returns the names in the archive in order.
func untar(t *testing.T, archive string, dir string) []string {
	t.Helper()
	f, err := os.Open(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()

	var names []string
	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, hdr.Name)
		target := filepath.Join(dir, filepath.FromSlash(hdr.Name))
		if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, data, 0640); err != nil {
			t.Fatal(err)
		}
	}
	return names
}

func TestWriteArchive(t *testing.T) {
	img, err := random.Image(1024, 3)
	if err != nil {
		t.Fatal(err)
	}
	wantDigest, err := img.Digest()
	if err != nil {
		t.Fatal(err)
	}
	wantConfig, err := img.ConfigName()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		ref         string
		wantName    string
		wantRefName string
		wantTags    []string
	}{
		{
			name:        "tag",
			ref:         "registry.k8s.io/kwok/kwok:v0.9.0",
			wantName:    "registry.k8s.io/kwok/kwok:v0.9.0",
			wantRefName: "registry.k8s.io/kwok/kwok:v0.9.0",
			wantTags:    []string{"registry.k8s.io/kwok/kwok:v0.9.0"},
		},
		{
			name:        "docker hub short name",
			ref:         "busybox",
			wantName:    "docker.io/library/busybox:latest",
			wantRefName: "docker.io/library/busybox:latest",
			wantTags:    []string{"docker.io/library/busybox:latest"},
		},
		{
			name:     "digest",
			ref:      "registry.k8s.io/kwok/kwok@" + wantDigest.String(),
			wantName: "registry.k8s.io/kwok/kwok@" + wantDigest.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := name.ParseReference(tt.ref)
			if err != nil {
				t.Fatal(err)
			}

			dir := t.TempDir()
			archive := filepath.Join(dir, "image.tar")
			if err := writeArchive(img, ref, archive); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(archive + ".tmp"); !os.IsNotExist(err) {
				t.Errorf("temporary file is left behind: %v", err)
			}

			layoutDir := filepath.Join(dir, "layout")
			names := untar(t, archive, layoutDir)
			seen := map[string]int{}
			for _, n := range names {
				seen[n]++
			}
			for _, n := range []string{"oci-layout", "index.json", "manifest.json", "blobs/sha256/" + wantDigest.Hex, "blobs/sha256/" + wantConfig.Hex} {
				if seen[n] != 1 {
					t.Errorf("expected exactly one %q in archive, got %d in %v", n, seen[n], names)
				}
			}

			// It is an OCI image layout
			index, err := layout.ImageIndexFromPath(layoutDir)
			if err != nil {
				t.Fatal(err)
			}
			im, err := index.IndexManifest()
			if err != nil {
				t.Fatal(err)
			}
			if len(im.Manifests) != 1 {
				t.Fatalf("expected one manifest in index, got %d", len(im.Manifests))
			}
			desc := im.Manifests[0]
			if desc.Digest != wantDigest {
				t.Errorf("index digest = %s, want %s", desc.Digest, wantDigest)
			}
			if got := desc.Annotations[annotationContainerdImageName]; got != tt.wantName {
				t.Errorf("%s = %q, want %q", annotationContainerdImageName, got, tt.wantName)
			}
			if got := desc.Annotations[annotationRefName]; got != tt.wantRefName {
				t.Errorf("%s = %q, want %q", annotationRefName, got, tt.wantRefName)
			}
			ociImg, err := index.Image(wantDigest)
			if err != nil {
				t.Fatal(err)
			}
			if err := validate.Image(ociImg); err != nil {
				t.Errorf("invalid OCI image: %v", err)
			}

			// It is a docker archive
			data, err := os.ReadFile(filepath.Join(layoutDir, "manifest.json"))
			if err != nil {
				t.Fatal(err)
			}
			var manifests []dockerManifest
			if err := json.Unmarshal(data, &manifests); err != nil {
				t.Fatal(err)
			}
			if len(manifests) != 1 {
				t.Fatalf("expected one entry in manifest.json, got %d", len(manifests))
			}
			if len(manifests[0].RepoTags) != len(tt.wantTags) || (len(tt.wantTags) == 1 && manifests[0].RepoTags[0] != tt.wantTags[0]) {
				t.Errorf("RepoTags = %v, want %v", manifests[0].RepoTags, tt.wantTags)
			}
			var tag *name.Tag
			if tg, ok := ref.(name.Tag); ok {
				tag = &tg
			}
			dockerImg, err := tarball.ImageFromPath(archive, tag)
			if err != nil {
				t.Fatal(err)
			}
			gotConfig, err := dockerImg.ConfigName()
			if err != nil {
				t.Fatal(err)
			}
			if gotConfig != wantConfig {
				t.Errorf("docker config = %s, want %s", gotConfig, wantConfig)
			}
			if err := validate.Image(dockerImg); err != nil {
				t.Errorf("invalid docker image: %v", err)
			}
		})
	}
}

func TestWriteArchiveSharedLayers(t *testing.T) {
	layer, err := random.Layer(1024, "application/vnd.oci.image.layer.v1.tar+gzip")
	if err != nil {
		t.Fatal(err)
	}
	digest, err := layer.Digest()
	if err != nil {
		t.Fatal(err)
	}
	img, err := appendLayers(layer, layer)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := name.ParseReference("localhost/kwok:test")
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	archive := filepath.Join(dir, "image.tar")
	if err := writeArchive(img, ref, archive); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, n := range untar(t, archive, filepath.Join(dir, "layout")) {
		if n == "blobs/sha256/"+digest.Hex {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected the shared layer once in archive, got %d", count)
	}
}

func appendLayers(layers ...containerregistryv1.Layer) (containerregistryv1.Image, error) {
	return mutate.AppendLayers(empty.Image, layers...)
}
