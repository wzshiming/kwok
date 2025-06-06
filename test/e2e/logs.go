/*
Copyright 2023 The Kubernetes Authors.

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

package e2e

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"sigs.k8s.io/kwok/pkg/utils/exec"
	"sigs.k8s.io/kwok/pkg/utils/path"
	"sigs.k8s.io/kwok/test/e2e/helper"
)

// CaseLogs creates a feature that tests logs
func CaseLogs(kwokctlPath, clusterName, nodeName, namespace, tmpDir string) *features.FeatureBuilder {
	node := helper.NewNodeBuilder(nodeName).
		Build()
	pod0 := helper.NewPodBuilder("pod0").
		WithNamespace(namespace).
		WithNodeName(nodeName).
		Build()

	return features.New("Pod Logs").
		Setup(helper.CreateNode(node)).
		Setup(helper.CreatePod(pod0)).
		Teardown(helper.DeletePod(pod0)).
		Teardown(helper.DeleteNode(node)).
		Assess("test logs", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("test logs")

			f, err := os.OpenFile(path.Join(tmpDir, "logs.log"), os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err = f.Close()
				if err != nil {
					t.Fatal(err)
				}
				_ = os.Remove(path.Join(tmpDir, "logs.log"))
			}()

			buf := bytes.NewBuffer(nil)

			cmd, err := exec.Command(exec.WithWriteTo(exec.WithFork(ctx, true), buf), kwokctlPath, "--name", clusterName, "kubectl", "logs", "-f", "-n", namespace, "pod0")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err = cmd.Process.Kill()
				if err != nil {
					t.Fatal(err)
				}
			}()

			_, _ = f.Write([]byte("2016-10-06T00:00:00Z stdout F log content 1\n"))
			_, _ = f.Write([]byte("2016-10-06T00:00:00Z stdout F log content 2\n"))
			_, _ = f.Write([]byte("2016-10-06T00:00:00Z stdout F log content 3\n"))

			want := "log content 1\nlog content 2\nlog content 3\n"
			for i := 0; i != 30; i++ {
				if buf.String() == want {
					break
				}
				time.Sleep(1 * time.Second)
			}
			if buf.String() != want {
				t.Fatalf("want %s, got %s", want, buf.String())
			}

			return ctx
		})
}
