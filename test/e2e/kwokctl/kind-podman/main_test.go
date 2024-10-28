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

// Package kind_podman_test is a test environment for kwok.
package kind_podman_test

import (
	"context"
	"flag"
	"os"
	"runtime"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/support/kwok"

	"sigs.k8s.io/kwok/pkg/consts"
	"sigs.k8s.io/kwok/pkg/log"
	"sigs.k8s.io/kwok/pkg/utils/exec"
	"sigs.k8s.io/kwok/pkg/utils/path"
	"sigs.k8s.io/kwok/test/e2e/helper"
)

var (
	runtimeEnv     = consts.RuntimeTypeKindPodman
	testEnv        env.Environment
	updateTestdata = false
	pwd            = os.Getenv("PWD")
	rootDir        = path.Join(pwd, "../../../..")
	logsDir        = path.Join(rootDir, "logs")
	clusterName    = envconf.RandomName("kwok-e2e-kind-podman", 24)
	namespace      = envconf.RandomName("ns", 16)
	testImage      = "localhost/kwok:test"
	kwokctlPath    = path.Join(rootDir, "bin", runtime.GOOS, runtime.GOARCH, "kwokctl"+helper.BinSuffix)
	baseArgs       = []string{
		"--kwok-controller-image=" + testImage,
		"--runtime=" + runtimeEnv,
		"--enable-metrics-server",
		"--wait=15m",
	}
)

func init() {
	_ = os.Setenv("KWOK_WORKDIR", path.Join(rootDir, "workdir"))
	flag.BoolVar(&updateTestdata, "update-testdata", false, "update all of testdata")
}

func TestMain(m *testing.M) {
	testEnv = helper.Environment()

	k := kwok.NewProvider().
		WithName(clusterName).
		WithPath(kwokctlPath)
	testEnv.Setup(
		helper.BuildKwokImage(rootDir, testImage, consts.RuntimeTypePodman),
		helper.BuildKwokctlBinary(rootDir),
		func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
			ctx, err := helper.CreateCluster(k, append(baseArgs,
				"--controller-port=10247",
				"--prometheus-port=9090",
				"--etcd-port=2400",
				"--kube-scheduler-port=10250",
				"--kube-controller-manager-port=10260",
				"--dashboard-port=6060",
				"--jaeger-port=16686",
				"--config="+path.Join(rootDir, "test/e2e"),
				"--config="+path.Join(rootDir, "kustomize/metrics/usage"),
			)...)(ctx, cfg)
			if err != nil {
				logger := log.FromContext(ctx)

				logger.Info("exporting logs for control plane")
				err := exec.Exec(ctx, "docker", "logs", "kwok-"+clusterName+"-control-plane")
				if err != nil {
					logger.ErrorContext(ctx, "failed to export logs for control plane", "err", err)
				}
			}
			return ctx, err
		},
		helper.CreateNamespace(namespace),
	)
	testEnv.Finish(
		helper.ExportLogs(k, logsDir),
		helper.DestroyCluster(k),
	)
	os.Exit(testEnv.Run(m))
}
