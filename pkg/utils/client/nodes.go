/*
Copyright 2024 The Kubernetes Authors.

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

package client

import (
	"fmt"
	"net/http"
	"context"
	"net"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"

	"sigs.k8s.io/kwok/pkg/log"
)

type TypedNodesClient interface {
	Nodes(name string) (kubernetes.Interface, error)
}

func (g *clientset) Nodes(name string) (kubernetes.Interface, error) {
	if cli, ok := g.nodeCache[name]; ok {
		return cli, nil
	}

	restConfig, err := g.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	newRestConfig := rest.CopyConfig(restConfig)

	// Make sure it doesn't reuse the connections
	newRestConfig.Proxy = http.ProxyFromEnvironment
	newRestConfig.Transport = nil
	var dialer net.Dialer
	newRestConfig.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		logger := log.FromContext(ctx)
		logger.Info("new connection", "address", address, "node", name)
		return dialer.DialContext(ctx, network, address)
	}

	cli, err := kubernetes.NewForConfig(newRestConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes dynamicClient: %w", err)
	}

	g.nodeCache[name] = cli
	return cli, nil
}
