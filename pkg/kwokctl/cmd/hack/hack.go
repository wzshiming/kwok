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

// Package hack defines a parent command for hack data
package hack

import (
	"context"
	"errors"
	"fmt"
	"os"

	"sigs.k8s.io/kwok/pkg/config"
	"sigs.k8s.io/kwok/pkg/kwokctl/runtime"
	"sigs.k8s.io/kwok/pkg/log"
	"sigs.k8s.io/kwok/pkg/utils/exec"
	"sigs.k8s.io/kwok/pkg/utils/path"

	"github.com/spf13/cobra"
)

type flagpole struct {
	Name string
}

// NewCommand returns a new cobra.Command for get
func NewCommand(ctx context.Context) *cobra.Command {
	flags := &flagpole{}

	cmd := &cobra.Command{
		Use:   "hack [command]",
		Short: "[experimental] Hack [get, put, del] resources in etcd without apiserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.Name = config.DefaultCluster
			err := runE(cmd.Context(), flags, args)
			if err != nil {
				return fmt.Errorf("%v: %w", args, err)
			}
			return nil
		},
	}
	cmd.DisableFlagParsing = true
	return cmd
}

func runE(ctx context.Context, flags *flagpole, args []string) error {
	name := config.ClusterName(flags.Name)
	workdir := path.Join(config.ClustersDir, flags.Name)

	logger := log.FromContext(ctx)
	logger = logger.With("cluster", flags.Name)
	ctx = log.NewContext(ctx, logger)

	rt, err := runtime.DefaultRegistry.Load(ctx, name, workdir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Warn("Cluster does not exist")
		}
		return err
	}

	err = rt.KectlInCluster(exec.WithStdIO(ctx), args...)

	if err != nil {
		return err
	}
	return nil
}
