package kwok

import (
	"context"
	"os"
	"testing"
	"time"

	"sigs.k8s.io/kwok/pkg/utils/path"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"

	"sigs.k8s.io/kwok/pkg/utils/exec"
)

var (
	testenv     env.Environment
	clusterName string
	namespace   string
	ctx, _      = context.WithTimeout(context.Background(), time.Hour)
	rootDir     = path.Join(os.Getenv("PWD"), "../../..")
)

func TestMain(m *testing.M) {
	testenv = env.New()
	clusterName = envconf.RandomName("kwok", 16)
	namespace = envconf.RandomName("kwok-ns", 16)

	// Build the kwok image
	image := "local/kwok:test"
	ctx := ctx
	ctx = exec.WithStdIO(ctx)
	err := exec.Exec(ctx, "make",
		"IMAGE_PREFIX=local",
		"VERSION=test",
		"-C", rootDir,
		"build-image")
	if err != nil {
		panic(err)
	}

	testenv.
		Setup(
			envfuncs.CreateKindCluster(clusterName),
			envfuncs.CreateNamespace(namespace),
			envfuncs.LoadDockerImageToCluster(clusterName, image),
			envfuncs.SetupCRDs(path.Join(rootDir, "kustomize/kwok"), ""),
		).
		Finish(
			envfuncs.DeleteNamespace(namespace),
			envfuncs.DestroyKindCluster(clusterName),
		)
	os.Exit(testenv.Run(m))
}
