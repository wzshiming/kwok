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

package kind

import (
	"reflect"
	"testing"
)

func Test_imageSaveArgs(t *testing.T) {
	tests := []struct {
		runtime string
		want    []string
	}{
		{
			runtime: "docker",
			want:    []string{"save", "image:tag", "-o", "/tmp/image.tar"},
		},
		{
			runtime: "podman",
			want:    []string{"save", "--format", "oci-archive", "image:tag", "-o", "/tmp/image.tar"},
		},
		{
			runtime: "nerdctl",
			want:    []string{"save", "image:tag", "-o", "/tmp/image.tar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.runtime, func(t *testing.T) {
			got := imageSaveArgs(tt.runtime, "image:tag", "/tmp/image.tar")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageSaveArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
