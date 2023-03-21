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

package clientsets

import (
	"context"
	"net"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/kwok/pkg/log"
	"sigs.k8s.io/kwok/pkg/utils/maps"
)

type clientsetManager struct {
	restConfig rest.Config
	clientsets maps.SyncMap[string, *client]
}

type client struct {
	Clientset *kubernetes.Clientset
	Connects  maps.SyncMap[*connExitHook, struct{}]
}

func NewClientsetManager(ctx context.Context, master, kubeconfig string) (Clientsets, error) {
	cfg, err := buildConfigFromFlags(ctx, master, kubeconfig)
	if err != nil {
		return nil, err
	}

	err = setConfigDefaults(cfg)
	if err != nil {
		return nil, err
	}

	return &clientsetManager{
		restConfig: *cfg,
		clientsets: maps.SyncMap[string, *client]{},
	}, nil
}

func (c *clientsetManager) Get(ctx context.Context, node string) (kubernetes.Interface, error) {
	if v, ok := c.clientsets.Load(node); ok {
		return v.Clientset, nil
	}

	logger := log.FromContext(ctx)
	var dialer net.Dialer
	dialFunc := func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := dialer.DialContext(ctx, network, address)
		if err != nil {
			return nil, err
		}

		clientset, ok := c.clientsets.Load(node)
		if !ok {
			logger.Warn("clientset not found",
				"node", node,
			)
		} else {
			connHoot := newConnExitHook(conn, clientset.Connects.Delete)
			clientset.Connects.Store(connHoot, struct{}{})
			conn = connHoot
		}

		return conn, nil
	}

	restConfig := c.restConfig
	restConfig.Dial = dialFunc
	clientset, err := kubernetes.NewForConfig(&restConfig)
	if err != nil {
		return nil, err
	}

	c.clientsets.Store(node, &client{
		Clientset: clientset,
		Connects:  maps.SyncMap[*connExitHook, struct{}]{},
	})
	return clientset, nil
}

func (c *clientsetManager) Close(ctx context.Context, node string) error {
	logger := log.FromContext(ctx)
	clientset, ok := c.clientsets.LoadAndDelete(node)
	if !ok {
		logger.Warn("clientset not found",
			"node", node,
		)
		return nil
	}
	clientset.Connects.Range(func(conn *connExitHook, value struct{}) bool {
		err := conn.Conn.Close()
		if err != nil {
			logger.Error("close conn", err)
		}
		return true
	})

	clientset.Clientset.DiscoveryClient.RESTClient().(*rest.RESTClient).Client.CloseIdleConnections()

	return nil
}
