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

package runtime

import (
	"context"
	"path/filepath"
	"testing"

	"sigs.k8s.io/kwok/pkg/apis/internalversion"
	"sigs.k8s.io/kwok/pkg/config"
)

func TestGetDefaultClusterExec(t *testing.T) {
	ctx := context.Background()

	clusterExec, err := getDefaultClusterExec()
	if err != nil {
		t.Fatal(err)
	}

	// The local exec must survive the round trip through the saved config,
	// kwok refuses to exec when it is unset.
	path := filepath.Join(t.TempDir(), ConfigName)
	err = config.Save(ctx, path, []config.InternalObject{clusterExec})
	if err != nil {
		t.Fatal(err)
	}
	objs, err := config.Load(ctx, path)
	if err != nil {
		t.Fatal(err)
	}

	got := config.FilterWithType[*internalversion.ClusterExec](objs)
	if len(got) != 1 {
		t.Fatalf("expected 1 ClusterExec, got %d", len(got))
	}
	if got[0].Spec.Selector != nil {
		t.Fatalf("expected to match all pods, got selector %+v", got[0].Spec.Selector)
	}
	execs := got[0].Spec.Execs
	if len(execs) != 1 || len(execs[0].Containers) != 0 || execs[0].Local == nil {
		t.Fatalf("expected a local exec for all containers, got %+v", execs)
	}
}

func TestGetDefaultClusterPortForward(t *testing.T) {
	pf, err := getDefaultClusterPortForward(0)
	if err != nil {
		t.Fatal(err)
	}
	if pf.Spec.Selector != nil {
		t.Fatalf("expected to match all pods, got selector %+v", pf.Spec.Selector)
	}
	forwards := pf.Spec.Forwards
	if len(forwards) != 1 || len(forwards[0].Ports) != 0 || forwards[0].Target == nil {
		t.Fatalf("expected a target for all ports, got %+v", forwards)
	}
	if forwards[0].Target.Port != 10247 {
		t.Fatalf("expected the default port of kwok, got %d", forwards[0].Target.Port)
	}

	pf, err = getDefaultClusterPortForward(12345)
	if err != nil {
		t.Fatal(err)
	}
	if got := pf.Spec.Forwards[0].Target.Port; got != 12345 {
		t.Fatalf("expected the overridden port, got %d", got)
	}
}
