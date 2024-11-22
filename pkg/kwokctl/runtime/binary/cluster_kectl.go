package binary

import (
	"context"

	"sigs.k8s.io/kwok/pkg/utils/format"
	"sigs.k8s.io/kwok/pkg/utils/net"
)

// KectlInCluster implements the kectl subcommand
func (c *Cluster) KectlInCluster(ctx context.Context, args ...string) error {
	config, err := c.Config(ctx)
	if err != nil {
		return err
	}
	conf := &config.Options
	return c.Kectl(ctx, append([]string{"--endpoints", net.LocalAddress + ":" + format.String(conf.EtcdPort)}, args...)...)
}
