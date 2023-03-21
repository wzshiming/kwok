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

	"k8s.io/client-go/kubernetes"
)

type clientsetOnce struct {
	kubernetes.Interface
}

func NewClientsetOnce(ctx context.Context, master, kubeconfig string) (Clientsets, error) {
	restConfig, err := buildConfigFromFlags(ctx, master, kubeconfig)
	if err != nil {
		return nil, err
	}

	err = setConfigDefaults(restConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &clientsetOnce{
		Interface: clientset,
	}, nil
}

func (c *clientsetOnce) Get(ctx context.Context, node string) (kubernetes.Interface, error) {
	return c.Interface, nil
}

func (c *clientsetOnce) Close(ctx context.Context, node string) error {
	return nil
}
