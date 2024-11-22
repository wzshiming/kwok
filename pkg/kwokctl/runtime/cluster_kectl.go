package runtime

import (
	"context"
)

func (c *Cluster) kectlPath(ctx context.Context) (string, error) {
	config, err := c.Config(ctx)
	if err != nil {
		return "", err
	}
	conf := &config.Options
	kectlPath, err := c.EnsureBinary(ctx, "kectl", conf.KedctlBinary)
	if err != nil {
		return "", err
	}
	return kectlPath, nil
}

func (c *Cluster) Kectl(ctx context.Context, args ...string) error {
	kectlPath, err := c.kectlPath(ctx)
	if err != nil {
		return err
	}

	return c.Exec(ctx, kectlPath, args...)
}
